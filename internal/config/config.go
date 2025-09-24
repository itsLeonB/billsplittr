package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const AppName = "Billsplittr"

type Config struct {
	App
	DB
	Google
	Valkey
}

type App struct {
	Env     string        `default:"debug"`
	Port    string        `default:"8080"`
	Timeout time.Duration `default:"10s"`
}

type Google struct {
	ServiceAccount string `split_words:"true" required:"true"`
	BillBucketName string `split_words:"true" required:"true"`
}

func Load() Config {
	var app App
	envconfig.MustProcess("APP", &app)

	var db DB
	envconfig.MustProcess("DB", &db)

	var google Google
	envconfig.MustProcess("GOOGLE", &google)

	var valkey Valkey
	envconfig.MustProcess("VALKEY", &valkey)

	return Config{
		App:    app,
		DB:     db,
		Google: google,
		Valkey: valkey,
	}
}
