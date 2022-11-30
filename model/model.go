package model

import (
	"gorm.io/gorm"

	"github.com/glebarez/sqlite"

	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/constant"
)

var (
	_db *gorm.DB

	tables = []any{
		&Client{},
		&Task{},
	}
)

func Init() error {
	dsn := config.GetString("db.dsn")
	if dsn == "" {
		dsn = "db.sqlite"
	}

	cfg := &gorm.Config{
		PrepareStmt: true,
	}
	if constant.Debug() {
		cfg.Logger = newLogger()
	}

	db, err := gorm.Open(sqlite.Open(dsn), cfg)
	if err != nil {
		return err
	}

	err = db.AutoMigrate(tables...)
	if err != nil {
		return err
	}

	log.Debug().Msgf("db migrate success")

	_db = db

	return nil
}
