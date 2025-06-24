package mapper

import (
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

func TransferMethodToResponse(transferMethod entity.TransferMethod) dto.TransferMethodResponse {
	return dto.TransferMethodResponse{
		ID:        transferMethod.ID,
		Name:      transferMethod.Name,
		Display:   transferMethod.Display,
		CreatedAt: transferMethod.CreatedAt,
		UpdatedAt: transferMethod.UpdatedAt,
		DeletedAt: transferMethod.DeletedAt.Time,
	}
}
