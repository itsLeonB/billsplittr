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
	*Services
}

func All(configs config.Config) *Provider {
	logger := ProvideLogger(config.AppName, configs.Env)
	dbs := ProvideDBs(logger, configs)
	repos := ProvideRepositories(dbs, configs.Google, logger)

	return &Provider{
		Logger:       logger,
		DBs:          dbs,
		Repositories: repos,
		Services:     ProvideServices(configs.Google, repos, logger),
	}
}

func (p *Provider) Shutdown() error {
	var err error
	if p.DBs != nil {
		if e := p.DBs.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
	}
	if p.Repositories != nil {
		if e := p.Repositories.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}
