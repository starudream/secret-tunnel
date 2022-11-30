package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"

	"github.com/starudream/go-lib/codec/json"
	"github.com/starudream/go-lib/errx"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/constant"
)

var e = binary.BigEndian

func Write(c net.Conn, msg any, sc ...int) error {
	b, ok := typeByteMap[reflect.TypeOf(msg)]
	if !ok {
		return fmt.Errorf("message: type error(%x)", b)
	}

	bs := json.MustMarshal(msg)
	bl := len(bs)

	if bl > constant.MessageSize {
		return fmt.Errorf("message: length greater than %d", constant.MessageSize)
	}

	bb := &bytes.Buffer{}
	_ = binary.Write(bb, e, b)
	_ = binary.Write(bb, e, uint16(bl))
	_ = binary.Write(bb, e, bs)

	_, err := c.Write(bb.Bytes())
	if err != nil {
		return err
	}

	if constant.Debug() {
		log.Debug().CallerSkipFrame(skip(sc...)).Msgf("write message: %x %s", b, bs)
	}

	return nil
}

func WriteL(c net.Conn, msg any, ls ...log.L) bool {
	err := Write(c, msg, skip())
	if err != nil {
		l := log.Logger()
		if len(ls) > 0 {
			l = ls[0]
		}
		l.Warn().CallerSkipFrame(1).Msgf("write message error: %v", err)
	}
	return err == nil
}

func ReadRaw(c net.Conn, sc ...int) (byte, []byte, error) {
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

	bs := make([]byte, bl)
	bn, err := io.ReadFull(c, bs)
	if err != nil {
		return 0, nil, err
	}
	if bn != int(bl) {
		return 0, nil, fmt.Errorf("message: format error")
	}

	if constant.Debug() {
		log.Debug().CallerSkipFrame(skip(sc...)).Msgf("read message: %x %s", b, bs)
	}

	return b, bs, nil
}

func Read(c net.Conn, sc ...int) (v any, err error) {
	b, bs, err := ReadRaw(c, skip(sc...))
	if err != nil {
		return nil, err
	}
	v = reflect.New(byteTypeMap[b]).Interface()
	return v, json.Unmarshal(bs, v)
}

func ReadL(c net.Conn, ls ...log.L) (v any, ne bool) {
	v, err := Read(c, skip())
	if err != nil && !errx.Is(err, io.EOF) {
		l := log.Logger()
		if len(ls) > 0 {
			l = ls[0]
		}
		l.Warn().CallerSkipFrame(1).Msgf("read message error: %v", err)
	}
	return v, err == nil
}

func skip(sc ...int) int {
	if len(sc) > 0 {
		return sc[0] + 1
	}
	return 1
}
