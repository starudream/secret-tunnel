package model

import (
	"database/sql"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/starudream/go-lib/errx"
)

var (
	ErrNoRows   = errors.New("record not found")
	ErrConflict = errors.New("record conflict")
)

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	if errx.Is(err, sql.ErrNoRows) || errx.Is(err, gorm.ErrRecordNotFound) {
		return ErrNoRows
	} else {
		es := err.Error()
		if strings.HasSuffix(es, "(1555)") || strings.HasSuffix(es, "(2067)") {
			return ErrConflict
		}
	}

	return err
}

// if se, ok := err.(*sqlite.Error); ok {
//		code := se.Code()
//		switch code {
//		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
//			return ErrConflict
//		}
//	}
