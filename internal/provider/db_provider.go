package provider

import (
	"errors"
	"fmt"

	"github.com/itsLeonB/billsplittr/internal/pkg/config"
	"github.com/itsLeonB/billsplittr/internal/pkg/logger"
	"github.com/itsLeonB/meq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBs struct {
	GormDB *gorm.DB
	MQ     meq.DB
}

func ProvideDBs() *DBs {
	dbs := &DBs{
		nil,
		meq.NewAsynqDB(logger.Global, config.Global.ToRedisOpts()),
	}

	dbs.openGormConnection()

	return dbs
}

func (d *DBs) Shutdown() error {
	var errs error

	if d.GormDB != nil {
		db, err := d.GormDB.DB()
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			if e := db.Close(); e != nil {
				errs = errors.Join(errs, e)
			}
		}
	}
	if d.MQ != nil {
		if e := d.MQ.Shutdown(); e != nil {
			errs = errors.Join(errs, e)
		}
	}

	return errs
}

func (d *DBs) Ping() error {
	var errs error

	if d.GormDB != nil {
		sqlDB, err := d.GormDB.DB()
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			if e := sqlDB.Ping(); e != nil {
				errs = errors.Join(errs, e)
			}
		}
	}
	if d.MQ != nil {
		if e := d.MQ.Ping(); e != nil {
			errs = errors.Join(errs, e)
		}
	}

	return errs
}

func (d *DBs) openGormConnection() {
	db, err := gorm.Open(postgres.Open(config.Global.ToPostgresDSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error opening GORM connection: %s", err.Error()))
	}

	d.GormDB = db
}
