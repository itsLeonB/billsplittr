package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App
	Auth
	DB
	Google
	ServiceClient
}

type App struct {
	Name       string        `default:"Cocoon"`
	Env        string        `default:"debug"`
	Port       string        `default:"50051"`
	Timeout    time.Duration `default:"10s"`
	ClientUrls []string      `split_words:"true"`
}

type Auth struct {
	SecretKey     string        `split_words:"true" default:"thisissecret"`
	TokenDuration time.Duration `split_words:"true" default:"24h"`
	Issuer        string        `default:"cocoon"`
	HashCost      int           `split_words:"true" default:"10"`
}

type Google struct {
	ServiceAccount string `split_words:"true" required:"true"`
}

type ServiceClient struct {
	CocoonHost string `split_words:"true" required:"true"`
}

func Load() Config {
	var app App
	envconfig.MustProcess("APP", &app)

	var auth Auth
	envconfig.MustProcess("AUTH", &auth)

	var db DB
	envconfig.MustProcess("DB", &db)

	var google Google
	envconfig.MustProcess("GOOGLE", &google)

	var svcClient ServiceClient
	envconfig.MustProcess("SERVICE_CLIENT", &svcClient)

	return Config{
		App:           app,
		Auth:          auth,
		DB:            db,
		Google:        google,
		ServiceClient: svcClient,
	}
}
