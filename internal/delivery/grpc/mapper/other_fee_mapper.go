package mapper

import (
	"github.com/itsLeonB/billsplittr-protos/gen/go/domain/v1"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/gerpc"
	"github.com/rotisserie/eris"
	"golang.org/x/text/currency"
)

func FromProtoFeeCalculationMethod(fcm domain.FeeCalculationMethod) (appconstant.FeeCalculationMethod, error) {
	switch fcm {
	case domain.FeeCalculationMethod_FEE_CALCULATION_METHOD_EQUAL_SPLIT:
		return appconstant.EqualSplitFee, nil
	case domain.FeeCalculationMethod_FEE_CALCULATION_METHOD_ITEMIZED_SPLIT:
		return appconstant.ItemizedSplitFee, nil
	default:
		return "", eris.Errorf("unknown fee calculation method enum: %s", fcm.String())
	}
}

func ToProtoFeeCalculationMethod(fcm appconstant.FeeCalculationMethod) (domain.FeeCalculationMethod, error) {
	switch fcm {
	case appconstant.EqualSplitFee:
		return domain.FeeCalculationMethod_FEE_CALCULATION_METHOD_EQUAL_SPLIT, nil
	case appconstant.ItemizedSplitFee:
		return domain.FeeCalculationMethod_FEE_CALCULATION_METHOD_ITEMIZED_SPLIT, nil
	default:
		return domain.FeeCalculationMethod_FEE_CALCULATION_METHOD_UNSPECIFIED, eris.Errorf("unknown fee calculation method constant: %s", fcm)
	}
}

func toFeeParticipantResponseProto(feeParticipant dto.FeeParticipantResponse) *domain.FeeParticipantResponse {
	return &domain.FeeParticipantResponse{
		ProfileId:   feeParticipant.ProfileID.String(),
		ShareAmount: ezutil.DecimalToMoney(feeParticipant.ShareAmount, currency.IDR.String()),
	}
}

func ToOtherFeeResponseProto(fee dto.OtherFeeResponse) (*domain.OtherFeeResponse, error) {
	calculationMethod, err := ToProtoFeeCalculationMethod(fee.CalculationMethod)
	if err != nil {
		return nil, err
	}
	return &domain.OtherFeeResponse{
		GroupExpenseId: fee.GroupExpenseID.String(),
		OtherFee: &domain.OtherFee{
			Name:              fee.Name,
			Amount:            ezutil.DecimalToMoney(fee.Amount, currency.IDR.String()),
			CalculationMethod: calculationMethod,
		},
		Participants: ezutil.MapSlice(fee.Participants, toFeeParticipantResponseProto),
		AuditMetadata: &domain.AuditMetadata{
			Id:        fee.ID.String(),
			CreatedAt: gerpc.NullableTimeToProto(fee.CreatedAt),
			UpdatedAt: gerpc.NullableTimeToProto(fee.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(fee.DeletedAt),
		},
	}, nil
}

func FromOtherFeeProto(fee *domain.OtherFee) (dto.OtherFeeData, error) {
	if fee == nil {
		return dto.OtherFeeData{}, eris.New("other fee input is nil")
	}

	calculationMethod, err := FromProtoFeeCalculationMethod(fee.GetCalculationMethod())
	if err != nil {
		return dto.OtherFeeData{}, err
	}

	return dto.OtherFeeData{
		Name:              fee.GetName(),
		Amount:            ezutil.MoneyToDecimal(fee.GetAmount()),
		CalculationMethod: calculationMethod,
	}, nil
}

func ToCalculationMethodInfoProto(method dto.FeeCalculationMethodInfo) (*domain.FeeCalculationMethodInfo, error) {
	enum, err := ToProtoFeeCalculationMethod(method.Method)
	if err != nil {
		return nil, err
	}
	return &domain.FeeCalculationMethodInfo{
		Method:      enum,
		Display:     method.Display,
		Description: method.Description,
	}, nil
}
