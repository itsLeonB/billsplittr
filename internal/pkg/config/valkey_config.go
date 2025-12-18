package config

import (
	"crypto/tls"

	"github.com/hibiken/asynq"
)

type Valkey struct {
	Addr      string
	Password  string
	Db        int
	EnableTls bool `split_words:"true" default:"true"`
}

func (v Valkey) ToRedisOpts() asynq.RedisClientOpt {
	var tlsConfig *tls.Config
	if v.EnableTls {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	return asynq.RedisClientOpt{
		Addr:      v.Addr,
		Password:  v.Password,
		DB:        v.Db,
		TLSConfig: tlsConfig,
	}
}
