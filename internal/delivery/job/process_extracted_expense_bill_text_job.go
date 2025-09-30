package job

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/provider"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/ezutil/v2"
)

type processExtractedExpenseBillTextJobImpl struct {
	groupExpenseService service.GroupExpenseService
}

func processExtractedExpenseBillTextJob(configs config.Config) (*ezutil.Job, error) {
	logger := provider.ProvideLogger("Process Extracted Expense Bill", configs.Env)

	providers, err := provider.All(configs, logger)
	if err != nil {
		return nil, err
	}

	jobImpl := processExtractedExpenseBillTextJobImpl{
		providers.Services.GroupExpense,
	}

	job := ezutil.NewJob(logger, jobImpl.Run).
		WithSetupFunc(providers.Ping).
		WithCleanupFunc(providers.Shutdown)

	return job, nil
}

func (j *processExtractedExpenseBillTextJobImpl) Run() error {
	return j.groupExpenseService.ParseFromBillText(context.Background())
}
