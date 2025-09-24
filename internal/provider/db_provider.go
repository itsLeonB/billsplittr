package provider

import (
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/itsLeonB/ezutil/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBs struct {
	dbConfig config.DB
	GormDB   *gorm.DB
	Asynq    *asynq.Client
}

func ProvideDBs(logger ezutil.Logger, cfg config.Config) *DBs {
	dbs := &DBs{
		cfg.DB,
		nil,
		connectAsynq(logger, cfg.Valkey),
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

	if d.Asynq != nil {
		if e := d.Asynq.Close(); e != nil {
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

	if d.Asynq != nil {
		if e := d.Asynq.Ping(); e != nil {
			errs = errors.Join(errs, e)
		}
	}

	return errs
}

func (d *DBs) getDSN() string {
	switch d.dbConfig.Driver {
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.dbConfig.User,
			d.dbConfig.Password,
			d.dbConfig.Host,
			d.dbConfig.Port,
			d.dbConfig.Name,
		)
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s",
			d.dbConfig.Host,
			d.dbConfig.User,
			d.dbConfig.Password,
			d.dbConfig.Name,
			d.dbConfig.Port,
		)
	default:
		panic(fmt.Sprintf("unsupported SQLDB driver: %s", d.dbConfig.Driver))
	}
}

func (d *DBs) getGormDialector() gorm.Dialector {
	switch d.dbConfig.Driver {
	// case "mysql":
	// 	return mysql.Open(sqldb.getDSN())
	case "postgres":
		return postgres.Open(d.getDSN())
	default:
		panic(fmt.Sprintf("unsupported SQLDB driver: %s", d.dbConfig.Driver))
	}
}

func (d *DBs) openGormConnection() {
	db, err := gorm.Open(d.getGormDialector(), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error opening GORM connection: %s", err.Error()))
	}

	d.GormDB = db
}

func connectAsynq(logger ezutil.Logger, cfg config.Valkey) *asynq.Client {
	if cfg.Addr == "" {
		logger.Warn("valkey config not provided, will not connect to asynq")
		return nil
	}

	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Db,
	})
}
