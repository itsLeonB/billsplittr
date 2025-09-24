package job

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
)

type cleanupOrphanedBillsJob struct {
	expenseBillSvc service.ExpenseBillService
}

func CleanupOrphanedBillsJob(configs config.Config) *ezutil.Job {
	var providers *provider.Provider
	jobImpl := cleanupOrphanedBillsJob{}
	logger := provider.ProvideLogger("Cleanup Orphaned Bills", configs.Env)

	return ezutil.NewJob(logger, jobImpl.Run).
		WithSetupFunc(func() error {
			providers = provider.All(configs)
			providers.Logger = logger
			jobImpl.expenseBillSvc = providers.Services.ExpenseBill
			return nil
		}).
		WithCleanupFunc(providers.Shutdown)
}

func (j *cleanupOrphanedBillsJob) Run() error {
	return j.expenseBillSvc.EnqueueCleanup(context.Background())
}
