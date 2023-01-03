package model

import (
	"database/sql"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrNoRows   = errors.New("record not found")
	ErrConflict = errors.New("record conflict")
)

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
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
