package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App
	DB
	Google
}

type App struct {
	Name    string        `default:"Billsplittr"`
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

	return Config{
		App:    app,
		DB:     db,
		Google: google,
	}
}
