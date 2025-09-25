package job

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
)

type enqueueCleanupOrphanedBillsJob struct {
	expenseBillSvc service.ExpenseBillService
}

func EnqueueCleanupOrphanedBillsJob(configs config.Config) *ezutil.Job {
	logger := provider.ProvideLogger("Enqueue Cleanup Orphaned Bills", configs.Env)
	providers := provider.All(configs, logger)
	jobImpl := enqueueCleanupOrphanedBillsJob{providers.Services.ExpenseBill}

	return ezutil.NewJob(logger, jobImpl.Run).
		WithSetupFunc(providers.Ping).
		WithCleanupFunc(providers.Shutdown)
}

func (j *enqueueCleanupOrphanedBillsJob) Run() error {
	return j.expenseBillSvc.EnqueueCleanup(context.Background())
}
