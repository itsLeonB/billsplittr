package service

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/mapper"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil"
)

type transferMethodServiceImpl struct {
	transferMethodRepository repository.TransferMethodRepository
}

func NewTransferMethodService(transferMethodRepository repository.TransferMethodRepository) TransferMethodService {
	return &transferMethodServiceImpl{transferMethodRepository}
}

func (tms *transferMethodServiceImpl) GetAll(ctx context.Context) ([]dto.TransferMethodResponse, error) {
	transferMethods, err := tms.transferMethodRepository.FindAll(ctx, ezutil.Specification[entity.TransferMethod]{})
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(transferMethods, mapper.TransferMethodToResponse), nil
}
