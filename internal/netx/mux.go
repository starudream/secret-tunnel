package netx

import (
	"bytes"
	"io"
	"time"

	"github.com/hashicorp/yamux"

	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/constant"
)

func yConfig() *yamux.Config {
	c := &yamux.Config{
		AcceptBacklog:          1024,
		EnableKeepAlive:        true,
		KeepAliveInterval:      30 * time.Second,
		ConnectionWriteTimeout: 10 * time.Second,
		MaxStreamWindowSize:    256 * 1024,
		StreamOpenTimeout:      5 * time.Second,
		StreamCloseTimeout:     60 * time.Second,
	}
	if constant.Debug() {
		c.LogOutput = &yLogger{l: log.With().Str("span", "mux").CallerWithSkipFrameCount(5).Logger()}
	} else {
		c.LogOutput = io.Discard
	}
	return c
}

type yLogger struct {
	l log.L
}

var _ io.Writer = (*yLogger)(nil)

var (
	yCE = []byte("[ERR]")
	yTE = []byte("[ERR] yamux: ")
	yCW = []byte("[WARN]")
	yTW = []byte("[WARN] yamux: ")
)

func (y *yLogger) Write(p []byte) (n int, err error) {
	if len(p) >= 20 {
		p = p[20:]
	}
	p = bytes.TrimSuffix(p, []byte("\n"))
	if bytes.Contains(p, yCE) {
		y.l.Error().Msg(string(bytes.TrimPrefix(p, yTE)))
	} else if bytes.Contains(p, yCW) {
		y.l.Warn().Msg(string(bytes.TrimPrefix(p, yTW)))
	} else {
		y.l.Debug().Msg(string(p))
	}
	return len(p), nil
}
