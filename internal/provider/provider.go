package provider

import (
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/ezutil/v2"
)

type Provider struct {
	Logger ezutil.Logger
	*DBs
	*Repositories
	*Services
}

func All(configs config.Config, logger ezutil.Logger) *Provider {
	dbs := ProvideDBs(logger, configs)
	repos := ProvideRepositories(dbs, logger)

	return &Provider{
		Logger:       logger,
		DBs:          dbs,
		Repositories: repos,
		Services:     ProvideServices(repos, logger, configs.Storage),
	}
}

func (p *Provider) Shutdown() error {
	if p.DBs != nil {
		return p.Shutdown()
	}
	return nil
}
