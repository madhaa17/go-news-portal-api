package config

import (
	"fmt"

	"news-app/database/seeds"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg *Config) ConnPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Psql.User, cfg.Psql.Password, cfg.Psql.Host, cfg.Psql.Port, cfg.Psql.DBName)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("[ConnPostgres-1] Failed to connect to database" + cfg.Psql.Host)
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[ConnPostgres-2] Failed to get database connection")
		return nil, err
	}

	seeds.SeedRoles(db)

	sqlDb.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	sqlDb.SetMaxIdleConns(cfg.Psql.DBMaxIdle)

	return &Postgres{DB: db}, nil
}
