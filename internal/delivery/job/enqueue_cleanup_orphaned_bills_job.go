package job

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/pkg/logger"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
)

type enqueueCleanupOrphanedBillsJobImpl struct {
	expenseBillSvc service.ExpenseBillService
}

func EnqueueCleanupOrphanedBillsJob() (*ezutil.Job, error) {
	providers, err := provider.All()
	if err != nil {
		return nil, err
	}

	jobImpl := enqueueCleanupOrphanedBillsJobImpl{providers.Services.ExpenseBill}

	job := ezutil.NewJob(logger.Global, jobImpl.Run).
		WithSetupFunc(providers.Ping).
		WithCleanupFunc(providers.Shutdown)

	return job, nil
}

func (j *enqueueCleanupOrphanedBillsJobImpl) Run() error {
	return j.expenseBillSvc.EnqueueCleanup(context.Background())
}
