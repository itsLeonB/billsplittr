package main

import (
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/delivery/job"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	j := job.EnqueueCleanupOrphanedBillsJob(config.Load())
	j.Run()
}
