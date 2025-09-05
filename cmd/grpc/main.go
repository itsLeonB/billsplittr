package main

import (
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/delivery/grpc"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	srv := grpc.Setup(config.Load())
	srv.Run()
}
