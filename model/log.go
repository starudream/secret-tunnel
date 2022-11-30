package model

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/starudream/go-lib/errx"
	"github.com/starudream/go-lib/log"
)

type iLogger struct {
	l log.L
}

const dbSlowThreshold = 500 * time.Millisecond

var _ logger.Interface = (*iLogger)(nil)

func newLogger() *iLogger {
	return &iLogger{
		l: log.With().Str("span", "db").CallerWithSkipFrameCount(5).Logger(),
	}
}

func (i *iLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return i
}

func (i *iLogger) Info(_ context.Context, s string, v ...any) {
	i.l.Info().Msgf(s, v...)
}

func (i *iLogger) Warn(_ context.Context, s string, v ...any) {
	i.l.Warn().Msgf(s, v...)
}

func (i *iLogger) Error(_ context.Context, s string, v ...any) {
	i.l.Error().Msgf(s, v...)
}

func (i *iLogger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)

	sql, rows := fc()

	l := func() log.L {
		x := i.l.With().Dur("took", elapsed)
		if rows >= 0 {
			x = x.Int64("rows", rows)
		}
		return x.Logger()
	}()

	if err != nil && !errx.Is(err, gorm.ErrRecordNotFound) {
		l.Error().Msgf("db error: %v -->> %s", err, sql)
	} else {
		if elapsed > dbSlowThreshold {
			l.Warn().Msg(sql)
		} else {
			l.Debug().Msg(sql)
		}
	}
}
