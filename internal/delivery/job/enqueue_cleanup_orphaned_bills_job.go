package job

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
)

type enqueueCleanupOrphanedBillsJobImpl struct {
	expenseBillSvc service.ExpenseBillService
}

func enqueueCleanupOrphanedBillsJob(configs config.Config) (*ezutil.Job, error) {
	logger := provider.ProvideLogger("Enqueue Cleanup Orphaned Bills", configs.Env)

	providers, err := provider.All(configs, logger)
	if err != nil {
		return nil, err
	}

	jobImpl := enqueueCleanupOrphanedBillsJobImpl{providers.Services.ExpenseBill}

	job := ezutil.NewJob(logger, jobImpl.Run).
		WithSetupFunc(providers.Ping).
		WithCleanupFunc(providers.Shutdown)

	return job, nil
}

func (j *enqueueCleanupOrphanedBillsJobImpl) Run() error {
	return j.expenseBillSvc.EnqueueCleanup(context.Background())
}
