package provider

import (
	"errors"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/ezutil/v2"
)

type Provider struct {
	Logger ezutil.Logger
	*DBs
	*Queues
	*Repositories
	*Services
}

func All(configs config.Config, logger ezutil.Logger) (*Provider, error) {
	dbs := ProvideDBs(logger, configs)
	repos := ProvideRepositories(dbs, logger)
	queues, err := ProvideQueues(logger, dbs.MQ)
	if err != nil {
		if e := dbs.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
		return nil, err
	}

	services, err := ProvideServices(repos, logger, configs, queues)
	if err != nil {
		if e := dbs.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
		return nil, err
	}

	return &Provider{
		Logger:       logger,
		DBs:          dbs,
		Queues:       queues,
		Repositories: repos,
		Services:     services,
	}, nil
}

func (p *Provider) Shutdown() error {
	if p.DBs != nil {
		return p.DBs.Shutdown()
	}
	return nil
}
