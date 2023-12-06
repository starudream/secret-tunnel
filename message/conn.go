package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"

	"github.com/starudream/go-lib/core/v2/codec/json"
	"github.com/starudream/go-lib/core/v2/slog"
)

type Conn struct {
	net.Conn
}

func NewConn(c net.Conn) *Conn {
	return &Conn{c}
}

var _ net.Conn = (*Conn)(nil)

func AcceptConn(ln net.Listener) (*Conn, error) {
	c, err := ln.Accept()
	if err != nil {
		return nil, err
	}
	return NewConn(c), nil
}

func (c *Conn) RemoteAddrString() string {
	if c == nil {
		return "<nil>"
	}
	return c.RemoteAddr().String()
}

func (c *Conn) SetReadTimeout(d ...time.Duration) {
	if len(d) > 0 {
		_ = c.SetReadDeadline(time.Now().Add(d[0]))
	} else {
		_ = c.SetReadDeadline(time.Time{})
	}
}

func (c *Conn) ReadMessage(d ...time.Duration) (any, error) {
	if len(d) > 0 {
		c.SetReadTimeout(d[0])
		defer c.SetReadTimeout()
	}
	return Read(c)
}

func (c *Conn) WriteMessage(msg any) bool {
	err := Write(c, msg)
	if err != nil {
		slog.Warn("write message error: %v", err)
	}
	return err == nil
}

var e = binary.BigEndian

func Write(c net.Conn, msg any) error {
	b, ok := typeByteMap[reflect.TypeOf(msg)]
	if !ok {
		return fmt.Errorf("message: type error(%x)", b)
	}

	bs := json.MustMarshal(msg)
	bl := len(bs)

	bb := &bytes.Buffer{}
	_ = binary.Write(bb, e, b)
	_ = binary.Write(bb, e, uint16(bl))
	_ = binary.Write(bb, e, bs)

	_, err := c.Write(bb.Bytes())
	if err != nil {
		return fmt.Errorf("message: write error: %w", err)
	}

	return nil
}

func ReadRaw(c net.Conn) (byte, []byte, error) {
	b := byte(0)
	err := binary.Read(c, e, &b)
	if err != nil {
		return 0, nil, err
	}

	_, ok := typeMap[b]
	if !ok {
		return 0, nil, fmt.Errorf("message: type error(%x)", b)
	}

	bl := uint16(0)
	err = binary.Read(c, e, &bl)
	if err != nil {
		return 0, nil, err
	}
	if bl > MaxSize {
		return 0, nil, fmt.Errorf("message: length greater than %d", MaxSize)
	}

	bs := make([]byte, bl)
	bn, err := io.ReadFull(c, bs)
	if err != nil {
		return 0, nil, err
	}
	if bn != int(bl) {
		return 0, nil, fmt.Errorf("message: format error")
	}

	return b, bs, nil
}

func Read(c net.Conn) (v any, err error) {
	b, bs, err := ReadRaw(c)
	if err != nil {
		return nil, err
	}
	v = reflect.New(byteTypeMap[b]).Interface()
	return v, json.Unmarshal(bs, v)
}
