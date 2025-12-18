package main

import (
	"github.com/itsLeonB/billsplittr/internal/delivery/job"
	"github.com/itsLeonB/billsplittr/internal/pkg/config"
	"github.com/itsLeonB/billsplittr/internal/pkg/logger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rotisserie/eris"
)

func main() {
	logger.Init()

	if err := config.Load(); err != nil {
		logger.Global.Fatal(eris.ToString(err, true))
	}

	j, err := job.EnqueueCleanupOrphanedBillsJob()
	if err != nil {
		logger.Global.Fatal(eris.ToString(err, true))
	}

	j.Run()
}
