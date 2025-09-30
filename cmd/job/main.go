package main

import (
	"log"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/delivery/job"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rotisserie/eris"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(eris.ToString(err, true))
	}
	j, err := job.SelectJob(cfg)
	if err != nil {
		log.Fatal(eris.ToString(err, true))
	}
	j.Run()
}
