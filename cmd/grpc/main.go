package main

import (
	"log"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rotisserie/eris"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(eris.ToString(err, true))
	}
	if cfg.ServiceType != appconstant.GRPCServer {
		log.Fatalf("wrong service type for grpc server entrypoint: %s", cfg.ServiceType)
	}
	srv, err := grpc.Setup(cfg)
	if err != nil {
		log.Fatal(eris.ToString(err, true))
	}
	srv.Run()
}
