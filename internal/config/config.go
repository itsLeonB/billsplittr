package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const AppName = "Billsplittr"

type Config struct {
	App
	DB
	Valkey
	Storage
}

type App struct {
	Env     string        `default:"debug"`
	Port    string        `default:"8080"`
	Timeout time.Duration `default:"10s"`
}

func Load() Config {
	var app App
	envconfig.MustProcess("APP", &app)

	var db DB
	envconfig.MustProcess("DB", &db)

	var valkey Valkey
	envconfig.MustProcess("VALKEY", &valkey)

	var storage Storage
	envconfig.MustProcess("STORAGE", &storage)

	return Config{
		App:     app,
		DB:      db,
		Valkey:  valkey,
		Storage: storage,
	}
}
