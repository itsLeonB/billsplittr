package job

import (
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
)

func SelectJob(cfg config.Config) (*ezutil.Job, error) {
	switch cfg.ServiceType {
	case appconstant.EnqueueCleanupOrphanedBillsJob:
		return enqueueCleanupOrphanedBillsJob(cfg)
	case appconstant.ProcessExtractedExpenseBillTextJob:
		return processExtractedExpenseBillJob(cfg)
	default:
		return nil, eris.Errorf("unknown service type for job entrypoint: %s", cfg.ServiceType)
	}
}
