package main

import (
	"github.com/itsLeonB/billsplittr/internal/delivery/http/server"
	"github.com/itsLeonB/billsplittr/internal/logging"
	"github.com/itsLeonB/ezutil"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	defaultConfigs := server.DefaultConfigs()
	logging.InitLogger(defaultConfigs.App)
	ezutil.RunServer(defaultConfigs, server.SetupHTTPServer)
}
