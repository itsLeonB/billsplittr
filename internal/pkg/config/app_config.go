package config

import "time"

const AppName = "Billsplittr"

type App struct {
	Env     string        `default:"debug"`
	Port    string        `default:"8080"`
	Timeout time.Duration `default:"10s"`
}
