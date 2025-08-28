package provider

import (
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/ezutil/v2"
)

type Provider struct {
	Logger ezutil.Logger
	*DBs
	*Repositories
	*Clients
	*Services
}

func All(configs config.Config) *Provider {
	dbs := ProvideDBs(configs.DB)
	repos := ProvideRepositories(dbs.GormDB, configs.Google)
	clients := ProvideClients(configs.ServiceClient)

	return &Provider{
		Logger:       ProvideLogger(configs.App),
		DBs:          dbs,
		Repositories: repos,
		Clients:      clients,
		Services:     ProvideServices(repos, clients),
	}
}
