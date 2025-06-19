package repository

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type friendshipRepositoryGorm struct {
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) FriendshipRepository {
	return &friendshipRepositoryGorm{db}
}

func (fr *friendshipRepositoryGorm) Insert(ctx context.Context, friendship entity.Friendship) (entity.Friendship, error) {
	db, err := fr.getGormInstance(ctx)
	if err != nil {
		return entity.Friendship{}, err
	}

	if err = db.Create(&friendship).Error; err != nil {
		return entity.Friendship{}, eris.Wrap(err, appconstant.MsgInsertData)
	}

	return friendship, nil
}

func (fr *friendshipRepositoryGorm) FindFirst(ctx context.Context, spec entity.FriendshipSpecification) (entity.Friendship, error) {
	var friendship entity.Friendship

	db, err := fr.getGormInstance(ctx)
	if err != nil {
		return entity.Friendship{}, err
	}

	err = db.
		Scopes(ezutil.WhereBySpec(spec.Friendship)).
		Joins("JOIN user_profiles AS up1 ON up1.id = friendships.profile_id_1").
		Joins("JOIN user_profiles AS up2 ON up2.id = friendships.profile_id_2").
		Where(
			fr.db.Where("up1.name = ? AND friendships.profile_id_1 <> ?", spec.Name, spec.ProfileID).
				Or("up2.name = ? AND friendships.profile_id_2 <> ?", spec.Name, spec.ProfileID),
		).
		Take(&friendship).
		Error

	if err != nil {
		return entity.Friendship{}, eris.Wrap(err, appconstant.MsgGetData)
	}

	return friendship, nil
}

func (fr *friendshipRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return fr.db.WithContext(ctx), nil
}
