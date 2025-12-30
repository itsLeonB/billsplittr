package provider

import (
	"errors"
)

type Provider struct {
	*DBs
	*Queues
	*Repositories
	*Services
}

func All() (*Provider, error) {
	dbs := ProvideDBs()
	repos := ProvideRepositories(dbs)
	queues := ProvideQueues(dbs.MQ)

	services, err := ProvideServices(repos, queues)
	if err != nil {
		if e := dbs.Shutdown(); e != nil {
			err = errors.Join(err, e)
		}
		return nil, err
	}

	return &Provider{
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
