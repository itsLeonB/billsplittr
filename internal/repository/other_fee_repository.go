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

type otherFeeRepositoryGorm struct {
	db *gorm.DB
	ezutil.CRUDRepository[entity.OtherFee]
}

func NewOtherFeeRepository(db *gorm.DB) OtherFeeRepository {
	return &otherFeeRepositoryGorm{
		db,
		ezutil.NewCRUDRepository[entity.OtherFee](db),
	}
}

func (ger *otherFeeRepositoryGorm) SyncParticipants(ctx context.Context, feeID uuid.UUID, participants []entity.FeeParticipant) error {
	db, err := ger.GetGormInstance(ctx)
	if err != nil {
		return err
	}

	profileIDs := make([]uuid.UUID, len(participants))
	for i, p := range participants {
		participants[i].OtherFeeID = feeID
		profileIDs[i] = p.ProfileID
	}

	if len(participants) > 0 {
		// For PostgreSQL
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "other_fee_id"}, {Name: "profile_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"share_amount"}),
		}).Create(&participants).Error; err != nil {
			return eris.Wrap(err, appconstant.ErrDataUpdate)
		}
	}

	query := db.Where("other_fee_id = ?", feeID)
	if len(profileIDs) > 0 {
		query = query.Where("profile_id NOT IN ?", profileIDs)
	}
	if err := query.Delete(&entity.ItemParticipant{}).Error; err != nil {
		return err
	}

	return nil
}
