package repository

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type userProfileRepositoryGorm struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepositoryGorm{db}
}

func (upr *userProfileRepositoryGorm) Insert(ctx context.Context, profile entity.UserProfile) (entity.UserProfile, error) {
	db, err := upr.getGormInstance(ctx)
	if err != nil {
		return entity.UserProfile{}, err
	}

	if err = db.Create(&profile).Error; err != nil {
		return entity.UserProfile{}, eris.Wrap(err, appconstant.MsgInsertData)
	}

	return profile, nil
}

func (upr *userProfileRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return upr.db.WithContext(ctx), nil
}
