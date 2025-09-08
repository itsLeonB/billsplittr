package mapper

import (
	"github.com/itsLeonB/audit/gen/go/audit/v1"
	"github.com/itsLeonB/billsplittr-protos/gen/go/otherfee/v1"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/gerpc"
	"github.com/rotisserie/eris"
	"golang.org/x/text/currency"
)

func FromProtoFeeCalculationMethod(fcm otherfee.FeeCalculationMethod) (appconstant.FeeCalculationMethod, error) {
	switch fcm {
	case otherfee.FeeCalculationMethod_FEE_CALCULATION_METHOD_EQUAL_SPLIT:
		return appconstant.EqualSplitFee, nil
	case otherfee.FeeCalculationMethod_FEE_CALCULATION_METHOD_ITEMIZED_SPLIT:
		return appconstant.ItemizedSplitFee, nil
	default:
		return "", eris.Errorf("unknown fee calculation method enum: %s", fcm.String())
	}
}

func ToProtoFeeCalculationMethod(fcm appconstant.FeeCalculationMethod) (otherfee.FeeCalculationMethod, error) {
	switch fcm {
	case appconstant.EqualSplitFee:
		return otherfee.FeeCalculationMethod_FEE_CALCULATION_METHOD_EQUAL_SPLIT, nil
	case appconstant.ItemizedSplitFee:
		return otherfee.FeeCalculationMethod_FEE_CALCULATION_METHOD_ITEMIZED_SPLIT, nil
	default:
		return otherfee.FeeCalculationMethod_FEE_CALCULATION_METHOD_UNSPECIFIED, eris.Errorf("unknown fee calculation method constant: %s", fcm)
	}
}

func toFeeParticipantResponseProto(feeParticipant dto.FeeParticipantResponse) *otherfee.FeeParticipantResponse {
	return &otherfee.FeeParticipantResponse{
		ProfileId:   feeParticipant.ProfileID.String(),
		ShareAmount: ezutil.DecimalToMoney(feeParticipant.ShareAmount, currency.IDR.String()),
	}
}

func ToOtherFeeResponseProto(fee dto.OtherFeeResponse) (*otherfee.OtherFeeResponse, error) {
	calculationMethod, err := ToProtoFeeCalculationMethod(fee.CalculationMethod)
	if err != nil {
		return nil, err
	}
	return &otherfee.OtherFeeResponse{
		GroupExpenseId: fee.GroupExpenseID.String(),
		OtherFee: &otherfee.OtherFee{
			Name:              fee.Name,
			Amount:            ezutil.DecimalToMoney(fee.Amount, currency.IDR.String()),
			CalculationMethod: calculationMethod,
		},
		Participants: ezutil.MapSlice(fee.Participants, toFeeParticipantResponseProto),
		AuditMetadata: &audit.Metadata{
			Id:        fee.ID.String(),
			CreatedAt: gerpc.NullableTimeToProto(fee.CreatedAt),
			UpdatedAt: gerpc.NullableTimeToProto(fee.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(fee.DeletedAt),
		},
	}, nil
}

func FromOtherFeeProto(fee *otherfee.OtherFee) (dto.OtherFeeData, error) {
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

func ToCalculationMethodInfoProto(method dto.FeeCalculationMethodInfo) (*otherfee.FeeCalculationMethodInfo, error) {
	enum, err := ToProtoFeeCalculationMethod(method.Method)
	if err != nil {
		return nil, err
	}
	return &otherfee.FeeCalculationMethodInfo{
		Method:      enum,
		Display:     method.Display,
		Description: method.Description,
	}, nil
}
