package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App
	DB
	Valkey
	Storage
	LLM
	Google
}

var Global Config

func Load() error {
	var err error

	var app App
	if e := envconfig.Process("APP", &app); e != nil {
		err = errors.Join(err, e)
	}

	var db DB
	if e := envconfig.Process("DB", &db); e != nil {
		err = errors.Join(err, e)
	}

	var valkey Valkey
	if e := envconfig.Process("VALKEY", &valkey); e != nil {
		err = errors.Join(err, e)
	}

	var storage Storage
	if e := envconfig.Process("STORAGE", &storage); e != nil {
		err = errors.Join(err, e)
	}

	var llm LLM
	if e := envconfig.Process("LLM", &llm); e != nil {
		err = errors.Join(err, e)
	}

	var google Google
	if e := envconfig.Process("GOOGLE", &google); e != nil {
		err = errors.Join(err, e)
	}

	if err != nil {
		return err
	}

	Global = Config{
		App:     app,
		DB:      db,
		Valkey:  valkey,
		Storage: storage,
		LLM:     llm,
		Google:  google,
	}

	return nil
}
