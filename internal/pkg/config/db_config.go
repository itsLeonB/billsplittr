package config

import "fmt"

type DB struct {
	Host     string `required:"true"`
	Port     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
	Name     string `required:"true" default:"billsplittr"`
}

func (d *DB) ToPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		d.Host,
		d.User,
		d.Password,
		d.Name,
		d.Port,
	)
}
