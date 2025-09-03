package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/drex-protos/gen/go/transaction/v1"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
)

type transferMethodServiceImpl struct {
	transferMethodClient transaction.TransferMethodServiceClient
}

func NewTransferMethodService(transferMethodClient transaction.TransferMethodServiceClient) TransferMethodService {
	return &transferMethodServiceImpl{transferMethodClient}
}

func (tms *transferMethodServiceImpl) GetAll(ctx context.Context) ([]dto.TransferMethodResponse, error) {
	response, err := tms.transferMethodClient.GetAll(ctx, nil)
	if err != nil {
		return nil, eris.Wrap(err, appconstant.ErrServiceClient)
	}

	mapFunc := func(res *transaction.TransferMethodResponse) (dto.TransferMethodResponse, error) {
		if res == nil {
			return dto.TransferMethodResponse{}, eris.New("nil transaction method response")
		}
		id, err := ezutil.Parse[uuid.UUID](res.GetId())
		if err != nil {
			return dto.TransferMethodResponse{}, err
		}

		return dto.TransferMethodResponse{
			ID:        id,
			Name:      res.GetName(),
			Display:   res.GetDisplay(),
			CreatedAt: ezutil.FromProtoTime(res.GetCreatedAt()),
			UpdatedAt: ezutil.FromProtoTime(res.GetUpdatedAt()),
			DeletedAt: ezutil.FromProtoTime(res.GetDeletedAt()),
		}, nil
	}

	return ezutil.MapSliceWithError(response.GetTransferMethods(), mapFunc)
}
