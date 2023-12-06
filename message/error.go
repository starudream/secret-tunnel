package message

import (
	"errors"
	"io"
	"net"
)

func ErrEOF(err error) bool {
	return err != nil && errors.Is(err, io.EOF)
}

func ErrClosed(err error) bool {
	return err != nil && errors.Is(err, net.ErrClosed)
}

func ErrOther(err error) bool {
	return err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed)
}
