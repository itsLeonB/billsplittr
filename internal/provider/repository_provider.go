package provider

import (
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil"
)

type Repositories struct {
	Transactor      ezutil.Transactor
	User            repository.UserRepository
	UserProfile     repository.UserProfileRepository
	Friendship      repository.FriendshipRepository
	DebtTransaction repository.DebtTransactionRepository
	TransferMethod  repository.TransferMethodRepository
}

func ProvideRepositories(configs *ezutil.Config) *Repositories {
	return &Repositories{
		Transactor:      ezutil.NewTransactor(configs.GORM),
		User:            repository.NewUserRepository(configs.GORM),
		UserProfile:     repository.NewUserProfileRepository(configs.GORM),
		Friendship:      repository.NewFriendshipRepository(configs.GORM),
		DebtTransaction: repository.NewDebtTransactionRepository(configs.GORM),
		TransferMethod:  repository.NewTransferMethodRepository(configs.GORM),
	}
}
