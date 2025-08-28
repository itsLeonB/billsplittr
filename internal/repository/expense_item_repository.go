package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	crud "github.com/itsLeonB/go-crud"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type expenseItemRepositoryGorm struct {
	db *gorm.DB
	crud.CRUDRepository[entity.ExpenseItem]
}

func NewExpenseItemRepository(db *gorm.DB) ExpenseItemRepository {
	return &expenseItemRepositoryGorm{
		db,
		crud.NewCRUDRepository[entity.ExpenseItem](db),
	}
}

func (ger *expenseItemRepositoryGorm) SyncParticipants(ctx context.Context, expenseItemID uuid.UUID, participants []entity.ItemParticipant) error {
	db, err := ger.GetGormInstance(ctx)
	if err != nil {
		return err
	}

	profileIDs := make([]uuid.UUID, len(participants))
	for i, p := range participants {
		participants[i].ExpenseItemID = expenseItemID
		profileIDs[i] = p.ProfileID
	}

	if len(participants) > 0 {
		// For PostgreSQL
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "expense_item_id"}, {Name: "profile_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"share"}),
		}).Create(&participants).Error; err != nil {
			return eris.Wrap(err, appconstant.ErrDataUpdate)
		}
	}

	query := db.Where("expense_item_id = ?", expenseItemID)
	if len(profileIDs) > 0 {
		query = query.Where("profile_id NOT IN ?", profileIDs)
	}
	if err := query.Delete(&entity.ItemParticipant{}).Error; err != nil {
		return err
	}

	return nil
}
