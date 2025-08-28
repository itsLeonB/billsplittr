package provider

import (
	"errors"

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

func (p *Provider) Shutdown() error {
	var err error
	if p.Clients != nil {
		if e := p.Clients.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
	}
	if p.DBs != nil {
		if e := p.DBs.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}
