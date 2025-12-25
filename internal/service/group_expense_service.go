package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/message"
	"github.com/itsLeonB/billsplittr/internal/pkg/logger"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/billsplittr/internal/service/fee"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type groupExpenseServiceImpl struct {
	transactor             crud.Transactor
	groupExpenseRepository repository.GroupExpenseRepository
	feeCalculatorRegistry  map[appconstant.FeeCalculationMethod]fee.FeeCalculator
	otherFeeRepository     repository.OtherFeeRepository
	expenseBillRepository  repository.ExpenseBillRepository
	llmService             LLMService
}

func NewGroupExpenseService(
	transactor crud.Transactor,
	groupExpenseRepository repository.GroupExpenseRepository,
	otherFeeRepository repository.OtherFeeRepository,
	expenseBillRepository repository.ExpenseBillRepository,
	llmService LLMService,
) GroupExpenseService {
	return &groupExpenseServiceImpl{
		transactor,
		groupExpenseRepository,
		fee.NewFeeCalculatorRegistry(),
		otherFeeRepository,
		expenseBillRepository,
		llmService,
	}
}

func (ges *groupExpenseServiceImpl) CreateDraft(ctx context.Context, request dto.NewGroupExpenseRequest) (dto.GroupExpenseResponse, error) {
	if err := ges.validate(request); err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	groupExpense := mapper.GroupExpenseRequestToEntity(request)

	insertedGroupExpense, err := ges.groupExpenseRepository.Insert(ctx, groupExpense)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return mapper.GroupExpenseToResponse(insertedGroupExpense), nil
}

func (ges *groupExpenseServiceImpl) UpdateDraft(ctx context.Context, expense entity.GroupExpense, request dto.NewGroupExpenseRequest) error {
	if err := ges.validate(request); err != nil {
		return err
	}

	mappedExpense := mapper.GroupExpenseRequestToEntity(request)
	expense.TotalAmount = mappedExpense.TotalAmount
	expense.ItemsTotal = mappedExpense.ItemsTotal
	expense.FeesTotal = mappedExpense.FeesTotal
	expense.Items = mappedExpense.Items
	expense.OtherFees = mappedExpense.OtherFees
	expense.Status = mappedExpense.Status

	if strings.HasPrefix(expense.Description, "Untitled Expense at") {
		expense.Description = request.Description
	}

	_, err := ges.groupExpenseRepository.Update(ctx, expense)
	return err
}

func (ges *groupExpenseServiceImpl) GetAllCreated(ctx context.Context, profileID uuid.UUID) ([]dto.GroupExpenseResponse, error) {
	spec := crud.Specification[entity.GroupExpense]{}
	spec.Model.CreatorProfileID = profileID
	spec.PreloadRelations = []string{"Items", "OtherFees"}

	groupExpenses, err := ges.groupExpenseRepository.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(groupExpenses, mapper.GroupExpenseToResponse), nil
}

func (ges *groupExpenseServiceImpl) GetDetails(ctx context.Context, id uuid.UUID) (dto.GroupExpenseResponse, error) {
	spec := crud.Specification[entity.GroupExpense]{}
	spec.Model.ID = id
	spec.PreloadRelations = []string{
		"Items",
		"OtherFees",
		"Items.Participants",
		"Participants",
		"Bill",
	}

	groupExpense, err := ges.getGroupExpense(ctx, spec)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return mapper.GroupExpenseToResponse(groupExpense), nil
}

func (ges *groupExpenseServiceImpl) ConfirmDraft(ctx context.Context, id, profileID uuid.UUID) (dto.GroupExpenseResponse, error) {
	var response dto.GroupExpenseResponse

	err := ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.GroupExpense]{}
		spec.Model.ID = id
		spec.PreloadRelations = []string{"Items", "OtherFees", "Items.Participants"}
		spec.ForUpdate = true

		groupExpense, err := ges.getGroupExpense(ctx, spec)
		if err != nil {
			return err
		}

		if groupExpense.Confirmed || groupExpense.Status == appconstant.ConfirmedExpense {
			return ungerr.UnprocessableEntityError("already confirmed")
		}

		if len(groupExpense.Items) < 1 {
			return ungerr.UnprocessableEntityError("cannot confirm empty items")
		}

		participantsByProfileID := make(map[uuid.UUID]*entity.ExpenseParticipant)
		for _, item := range groupExpense.Items {
			if len(item.Participants) < 1 {
				return ungerr.UnprocessableEntityError(fmt.Sprintf("item %s does not have participants", item.Name))
			}
			for _, participant := range item.Participants {
				amountToAdd := item.TotalAmount().Mul(participant.Share)
				if expenseParticipant, ok := participantsByProfileID[participant.ProfileID]; ok {
					expenseParticipant.ShareAmount = expenseParticipant.ShareAmount.Add(amountToAdd)
				} else {
					expenseParticipant := entity.ExpenseParticipant{
						ParticipantProfileID: participant.ProfileID,
						ShareAmount:          amountToAdd,
					}
					participantsByProfileID[participant.ProfileID] = &expenseParticipant
				}
			}
		}

		groupExpenseParticipants := make([]entity.ExpenseParticipant, 0, len(participantsByProfileID))
		for _, expenseParticipant := range participantsByProfileID {
			groupExpenseParticipants = append(groupExpenseParticipants, *expenseParticipant)
		}

		groupExpense.Participants = groupExpenseParticipants
		updatedOtherFees, err := ges.calculateOtherFeeSplits(ctx, groupExpense)
		if err != nil {
			return err
		}

		for _, fee := range updatedOtherFees {
			for _, participant := range fee.Participants {
				if expenseParticipant, ok := participantsByProfileID[participant.ProfileID]; !ok {
					return eris.New("missing participant profile from other fee")
				} else {
					expenseParticipant.ShareAmount = expenseParticipant.ShareAmount.Add(participant.ShareAmount)
				}
			}
		}

		updatedGroupExpenseParticipants := make([]entity.ExpenseParticipant, 0, len(participantsByProfileID))
		for _, expenseParticipant := range participantsByProfileID {
			updatedGroupExpenseParticipants = append(updatedGroupExpenseParticipants, *expenseParticipant)
		}

		if err = ges.groupExpenseRepository.SyncParticipants(ctx, groupExpense.ID, updatedGroupExpenseParticipants); err != nil {
			return err
		}

		groupExpense.Confirmed = true
		// TODO: explore cleaner way
		groupExpense.Participants = nil // Prevent GORM updating child, already synced above

		updatedGroupExpense, err := ges.groupExpenseRepository.Update(ctx, groupExpense)
		if err != nil {
			return err
		}

		updatedGroupExpense.Participants = updatedGroupExpenseParticipants

		response = mapper.GroupExpenseToResponse(updatedGroupExpense)

		return nil
	})

	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return response, nil
}

func (ges *groupExpenseServiceImpl) getGroupExpense(ctx context.Context, spec crud.Specification[entity.GroupExpense]) (entity.GroupExpense, error) {
	groupExpense, err := ges.groupExpenseRepository.FindFirst(ctx, spec)
	if err != nil {
		return entity.GroupExpense{}, err
	}
	if groupExpense.IsZero() {
		return entity.GroupExpense{}, ungerr.NotFoundError(fmt.Sprintf("group expense with ID %s is not found", spec.Model.ID))
	}
	if groupExpense.IsDeleted() {
		return entity.GroupExpense{}, ungerr.UnprocessableEntityError(fmt.Sprintf("group expense with ID %s is deleted", spec.Model.ID))
	}

	return groupExpense, nil
}

func (ges *groupExpenseServiceImpl) validate(request dto.NewGroupExpenseRequest) error {
	if request.TotalAmount.IsZero() {
		return ungerr.UnprocessableEntityError(appconstant.ErrAmountZero)
	}

	calculatedFeeTotal := decimal.Zero
	calculatedSubtotal := decimal.Zero
	for _, item := range request.Items {
		calculatedSubtotal = calculatedSubtotal.Add(item.Amount.Mul(decimal.NewFromInt(int64(item.Quantity))))
	}
	for _, fee := range request.OtherFees {
		calculatedFeeTotal = calculatedFeeTotal.Add(fee.Amount)
	}
	if calculatedFeeTotal.Add(calculatedSubtotal).Cmp(request.TotalAmount) != 0 {
		return ungerr.UnprocessableEntityError(appconstant.ErrAmountMismatched)
	}
	if calculatedSubtotal.Cmp(request.Subtotal) != 0 {
		return ungerr.UnprocessableEntityError(appconstant.ErrAmountMismatched)
	}

	return nil
}

func (ges *groupExpenseServiceImpl) calculateOtherFeeSplits(ctx context.Context, groupExpense entity.GroupExpense) ([]entity.OtherFee, error) {
	var splitErr error

	mapperFunc := func(fee entity.OtherFee) entity.OtherFee {
		feeCalculator, ok := ges.feeCalculatorRegistry[fee.CalculationMethod]
		if !ok {
			splitErr = eris.Errorf("unsupported calculation method: %s", fee.CalculationMethod)
			return entity.OtherFee{}
		}

		if err := feeCalculator.Validate(fee, groupExpense); err != nil {
			splitErr = err
			return entity.OtherFee{}
		}

		fee.Participants = feeCalculator.Split(fee, groupExpense)

		if err := ges.otherFeeRepository.SyncParticipants(ctx, fee.ID, fee.Participants); err != nil {
			splitErr = err
			return entity.OtherFee{}
		}

		return fee
	}

	splitFees := ezutil.MapSlice(groupExpense.OtherFees, mapperFunc)

	return splitFees, splitErr
}

func (ges *groupExpenseServiceImpl) GetUnconfirmedGroupExpenseForUpdate(ctx context.Context, profileID, id uuid.UUID) (entity.GroupExpense, error) {
	spec := crud.Specification[entity.GroupExpense]{}
	spec.Model.ID = id
	spec.Model.CreatorProfileID = profileID
	spec.ForUpdate = true
	spec.PreloadRelations = []string{"Items"}
	groupExpense, err := ges.getGroupExpense(ctx, spec)
	if err != nil {
		return entity.GroupExpense{}, err
	}
	if groupExpense.Confirmed || groupExpense.Status == appconstant.ConfirmedExpense {
		return entity.GroupExpense{}, ungerr.UnprocessableEntityError("expense already confirmed")
	}

	return groupExpense, nil
}

func (ges *groupExpenseServiceImpl) ParseFromBillText(ctx context.Context, msg message.ExpenseBillTextExtracted) error {
	return ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		expenseBill, err := ges.getPendingForProcessingExpenseBill(ctx, msg.ID)
		if err != nil {
			return err
		}
		if expenseBill.IsDeleted() {
			logger.Global.Infof("expense bill ID: %s is deleted, skipping", expenseBill.ID.String())
			return nil
		}
		return ges.executeParseFlow(ctx, msg.Text, expenseBill)
	})
}

func (ges *groupExpenseServiceImpl) executeParseFlow(ctx context.Context, text string, expenseBill entity.ExpenseBill) error {
	if !expenseBill.GroupExpenseID.Valid {
		// No groupExpenseID yet, execute flow v1
		return ges.parseFlowV1(ctx, text, expenseBill)
	}

	if err := ges.parseFlowV2(ctx, text, expenseBill); err != nil {
		// Log error but do not return error (commit the transaction)
		logger.Global.Errorf("error processing bill parsing v2: %s", eris.ToString(err, true))
		return nil
	}

	return nil
}

func (ges *groupExpenseServiceImpl) parseFlowV1(ctx context.Context, text string, expenseBill entity.ExpenseBill) error {
	request, err := ges.parseExpenseBillTextToExpenseRequest(ctx, text)
	if err != nil {
		if errors.Is(err, appconstant.ErrExpenseNotDetected) {
			return nil
		}
		return err
	}

	request.CreatorProfileID = expenseBill.CreatorProfileID
	request.PayerProfileID = expenseBill.PayerProfileID

	groupExpense, err := ges.CreateDraft(ctx, request)
	if err != nil {
		return err
	}

	expenseBill.GroupExpenseID = uuid.NullUUID{
		UUID:  groupExpense.ID,
		Valid: true,
	}

	_, err = ges.expenseBillRepository.Update(ctx, expenseBill)

	return err
}

func (ges *groupExpenseServiceImpl) parseFlowV2(ctx context.Context, text string, expenseBill entity.ExpenseBill) error {
	status, err := ges.processAndGetStatus(ctx, text, expenseBill)
	expenseBill.Status = status
	_, statusErr := ges.expenseBillRepository.Update(ctx, expenseBill)
	if statusErr != nil {
		return errors.Join(statusErr, err)
	}
	return err
}

func (ges *groupExpenseServiceImpl) processAndGetStatus(ctx context.Context, text string, expenseBill entity.ExpenseBill) (appconstant.BillStatus, error) {
	expense, err := ges.GetUnconfirmedGroupExpenseForUpdate(ctx, expenseBill.CreatorProfileID, expenseBill.GroupExpenseID.UUID)
	if err != nil {
		return appconstant.FailedParsingBill, err
	}

	request, err := ges.parseExpenseBillTextToExpenseRequest(ctx, text)
	if err != nil {
		if errors.Is(err, appconstant.ErrExpenseNotDetected) {
			return appconstant.NotDetectedBill, nil
		}
		return appconstant.FailedParsingBill, err
	}

	if err = ges.UpdateDraft(ctx, expense, request); err != nil {
		return appconstant.FailedParsingBill, err
	}

	return appconstant.ParsedBill, nil
}

func (ges *groupExpenseServiceImpl) buildSystemPrompt() string {
	return `You are an expert at parsing expense documents and receipts. 
Extract the expense information and return ONLY a valid JSON object in the following schema:

{
  "totalAmount": number,
  "subtotal": number,
  "description": string,
  "items": [
    {
      "name": string,
      "amount": number,   // price per unit
      "quantity": number
    }
  ],
  "otherFees": [
    {
      "name": string,
      "amount": number,
      "calculationMethod": "EQUAL_SPLIT" | "ITEMIZED_SPLIT"
    }
  ]
}

INSTRUCTIONS:
1. Return ONLY the JSON object, no explanations, no comments.
2. The JSON must be compact (no spaces, no line breaks, no pretty formatting).
3. totalAmount = subtotal + sum of otherFees.
4. subtotal = sum of (item.amount * item.quantity).
5. If subtotal is not explicitly mentioned, calculate it.
6. If quantity is not specified, assume 1.
7. Item.amount is always price per unit, not total for all units.
8. For otherFees:
   - Use "ITEMIZED_SPLIT" for percentage-based fees like tax or service charge, 
     because they should be distributed proportionally to the items each person ordered.
   - Use "EQUAL_SPLIT" only for true flat fees that apply equally regardless of items 
     (e.g., table charge, fixed booking fee).
9. All numeric values must be numbers, not strings.
10. Decimal separator can be "." or ",". Normalize both to "." in the output.
   - Example: "10,5" → 10.5
   - Example: "10.50" → 10.5
11. Thousands separators may appear as "." or "," in the input. Always remove them before parsing.
   - Example: "10.000" → 10000
   - Example: "10,000" → 10000
12. The final output must contain plain numeric values, with no thousands separators, and "." as the decimal separator.
13. If no clear expense information exists, return string "NOT_DETECTED"`
}

func (ges *groupExpenseServiceImpl) buildUserPrompt(text string) string {
	return fmt.Sprintf("TEXT TO PARSE:\n%s", text)
}

func (ges *groupExpenseServiceImpl) getPendingForProcessingExpenseBill(ctx context.Context, id uuid.UUID) (entity.ExpenseBill, error) {
	spec := crud.Specification[entity.ExpenseBill]{}
	spec.Model.ID = id
	spec.ForUpdate = true
	spec.DeletedFilter = crud.IncludeDeleted
	expenseBill, err := ges.expenseBillRepository.FindFirst(ctx, spec)
	if err != nil {
		return entity.ExpenseBill{}, err
	}
	if expenseBill.Status == appconstant.ParsedBill {
		return entity.ExpenseBill{}, eris.Errorf("expense bill ID: %s already parsed", expenseBill.GroupExpenseID.UUID.String())
	}

	return expenseBill, nil
}

func (ges *groupExpenseServiceImpl) parseExpenseBillTextToExpenseRequest(ctx context.Context, text string) (dto.NewGroupExpenseRequest, error) {
	promptResponse, err := ges.llmService.Prompt(ctx, ges.buildSystemPrompt(), ges.buildUserPrompt(text))
	if err != nil {
		return dto.NewGroupExpenseRequest{}, err
	}
	if promptResponse == "NOT_DETECTED" {
		logger.Global.Info("group expense not detected")
		return dto.NewGroupExpenseRequest{}, appconstant.ErrExpenseNotDetected
	}

	var request dto.NewGroupExpenseRequest
	if err = json.Unmarshal([]byte(promptResponse), &request); err != nil {
		return dto.NewGroupExpenseRequest{}, eris.Wrap(err, "error unmarshaling to JSON")
	}

	return request, nil
}

// region V2

func (ges *groupExpenseServiceImpl) CreateDraftV2(ctx context.Context, req dto.NewDraftExpense) (dto.GroupExpenseResponse, error) {
	description := req.Description
	if req.Description == "" {
		description = "Untitled Expense at " + time.Now().Format(time.DateOnly)
	}

	newDraftExpense := entity.GroupExpense{
		CreatorProfileID: req.CreatorProfileID,
		Description:      description,
		Status:           appconstant.DraftExpense,
	}

	insertedDraftExpense, err := ges.groupExpenseRepository.Insert(ctx, newDraftExpense)
	if err != nil {
		return dto.GroupExpenseResponse{}, err
	}

	return mapper.GroupExpenseToResponse(insertedDraftExpense), nil
}

func (ges *groupExpenseServiceImpl) Delete(ctx context.Context, id, profileID uuid.UUID) error {
	return ges.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		expense, err := ges.GetUnconfirmedGroupExpenseForUpdate(ctx, profileID, id)
		if err != nil {
			return err
		}

		return ges.groupExpenseRepository.Delete(ctx, expense)
	})
}
