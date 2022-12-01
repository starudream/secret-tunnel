package model

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite"
)

func TestMain(m *testing.M) {
	err := InitTest()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
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

	fmt.Println("db migrate success")

	_db = db

	return nil
}
