package netx

import (
	"net"
	"time"

	"github.com/hashicorp/yamux"

	"github.com/starudream/secret-tunnel/constant"
)

type Conn struct {
	conn *net.TCPConn

	session *yamux.Session
}

func New(c net.Conn, isServer bool) *Conn {
	tc := c.(*net.TCPConn)
	_ = tc.SetReadBuffer(constant.MessageSize)
	_ = tc.SetWriteBuffer(constant.MessageSize)
	_ = tc.SetKeepAlive(true)
	_ = tc.SetKeepAlivePeriod(30 * time.Second)
	_ = tc.SetNoDelay(true)

	session := func() *yamux.Session {
		if isServer {
			session, _ := yamux.Server(tc, yConfig())
			return session
		}
		session, _ := yamux.Client(tc, yConfig())
		return session
	}()

	return &Conn{conn: tc, session: session}
}

func (c *Conn) Session() *yamux.Session {
	return c.session
}

func (c *Conn) Close() {
	if c == nil {
		return
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
	if c.session != nil {
		_ = c.session.Close()
	}
}

func RemoteAddrString(c net.Conn) string {
	if c == nil {
		return ""
	}
	return c.RemoteAddr().String()
}

func SetReadTimeout(conn net.Conn, tds ...time.Duration) {
	if len(tds) > 0 && tds[0] > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(tds[0]))
	} else {
		_ = conn.SetReadDeadline(time.Time{})
	}
}

func Close(conn net.Conn) {
	if conn == nil {
		return
	}
	_ = conn.Close()
}
