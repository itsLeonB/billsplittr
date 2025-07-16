package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type groupExpenseRepositoryGorm struct {
	ezutil.CRUDRepository[entity.GroupExpense]
	db *gorm.DB
}

func NewGroupExpenseRepository(db *gorm.DB) GroupExpenseRepository {
	return &groupExpenseRepositoryGorm{
		ezutil.NewCRUDRepository[entity.GroupExpense](db),
		db,
	}
}

func (ger *groupExpenseRepositoryGorm) SyncParticipants(ctx context.Context, groupExpenseID uuid.UUID, participants []entity.ExpenseParticipant) error {
	db, err := ger.GetGormInstance(ctx)
	if err != nil {
		return err
	}

	profileIDs := make([]uuid.UUID, len(participants))
	for i, p := range participants {
		participants[i].GroupExpenseID = groupExpenseID
		profileIDs[i] = p.ParticipantProfileID
	}

	if len(participants) > 0 {
		// For PostgreSQL
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "group_expense_id"}, {Name: "participant_profile_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"share_amount"}),
		}).Create(&participants).Error; err != nil {
			return eris.Wrap(err, appconstant.ErrDataUpdate)
		}
	}

	query := db.Where("group_expense_id = ?", groupExpenseID)
	if len(profileIDs) > 0 {
		query = query.Where("participant_profile_id NOT IN ?", profileIDs)
	}
	if err := query.Delete(&entity.ExpenseParticipant{}).Error; err != nil {
		return err
	}

	return nil
}
