package message

import (
	"io"
	"net"
	"time"

	"github.com/hashicorp/yamux"

	"github.com/starudream/go-lib/core/v2/gh"
)

type MuxConn struct {
	conn *net.TCPConn
	sess *yamux.Session
}

func NewMuxConn(c net.Conn, isServer bool) *MuxConn {
	tc := c.(*net.TCPConn)

	_ = tc.SetReadBuffer(MaxSize)
	_ = tc.SetWriteBuffer(MaxSize)
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(30 * time.Second)
	_ = tc.SetNoDelay(true)

	sess := func() *yamux.Session {
		cfg := &yamux.Config{
			AcceptBacklog:          1024,
			EnableKeepAlive:        true,
			KeepAliveInterval:      30 * time.Second,
			ConnectionWriteTimeout: 10 * time.Second,
			MaxStreamWindowSize:    256 * 1024,
			StreamOpenTimeout:      5 * time.Second,
			StreamCloseTimeout:     60 * time.Second,
			LogOutput:              io.Discard,
		}
		if isServer {
			session, _ := yamux.Server(tc, cfg)
			return session
		}
		session, _ := yamux.Client(tc, cfg)
		return session
	}()

	return &MuxConn{conn: tc, sess: sess}
}

func (c *MuxConn) Session() *yamux.Session {
	return c.sess
}

func (c *MuxConn) Open() (*Conn, error) {
	conn, err := c.sess.Open()
	if err != nil {
		return nil, err
	}
	return NewConn(conn), nil
}

func (c *MuxConn) Close() {
	if c == nil {
		return
	}
	if c.sess != nil {
		gh.Close(c.sess)
	}
	if c.conn != nil {
		gh.Close(c.conn)
	}
}
