package logger

import (
	"os"

	"github.com/itsLeonB/ezutil/v2"
)

var Global ezutil.Logger

func Init() {
	logLevel := 1
	level := os.Getenv("LOG_LEVEL")
	if level == "debug" {
		logLevel = 0
	}
	Global = ezutil.NewSimpleLogger("Billsplittr", true, logLevel)
}
