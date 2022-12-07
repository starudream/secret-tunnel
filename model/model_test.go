package model

import (
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"

	"github.com/starudream/go-lib/log"
)

func TestMain(m *testing.M) {
	err := InitTest()
	if err != nil {
		log.Fatal().Msgf("init error: %v", err.Error())
	}
	os.Exit(m.Run())
}

func InitTest() error {
	cfg := &gorm.Config{
		PrepareStmt: true,
		Logger:      newLogger(),
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), cfg)
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
