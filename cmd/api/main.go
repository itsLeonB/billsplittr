package main

import (
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/billsplittr/internal/delivery/http"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	srv := http.Setup(config.Load())
	srv.ServeGracefully()
}
