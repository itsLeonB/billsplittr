package provider

import (
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/itsLeonB/ezutil"
)

type Repositories struct {
	User repository.UserRepository
}

func ProvideRepositories(configs *ezutil.Config) *Repositories {
	return &Repositories{
		User: repository.NewUserRepository(configs.GORM),
	}
}
