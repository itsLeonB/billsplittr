package mapper

import (
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil/v2"
)

func OtherFeeToResponse(fee entity.OtherFee) dto.OtherFeeResponse {
	return dto.OtherFeeResponse{
		ID:                fee.ID,
		GroupExpenseID:    fee.GroupExpenseID,
		Name:              fee.Name,
		Amount:            fee.Amount,
		CalculationMethod: fee.CalculationMethod,
		CreatedAt:         fee.CreatedAt,
		UpdatedAt:         fee.UpdatedAt,
		DeletedAt:         fee.DeletedAt.Time,
		Participants:      ezutil.MapSlice(fee.Participants, feeParticipantToResponse),
	}
}

func feeParticipantToResponse(feeParticipant entity.FeeParticipant) dto.FeeParticipantResponse {
	return dto.FeeParticipantResponse{
		ProfileID:   feeParticipant.ProfileID,
		ShareAmount: feeParticipant.ShareAmount,
	}
}

func OtherFeeRequestToEntity(request dto.OtherFeeData) entity.OtherFee {
	return entity.OtherFee{
		Name:              request.Name,
		Amount:            request.Amount,
		CalculationMethod: request.CalculationMethod,
	}
}

func PatchOtherFeeWithRequest(otherFee entity.OtherFee, request dto.UpdateOtherFeeRequest) entity.OtherFee {
	otherFee.Name = request.Name
	otherFee.Amount = request.Amount
	otherFee.CalculationMethod = request.CalculationMethod
	return otherFee
}
