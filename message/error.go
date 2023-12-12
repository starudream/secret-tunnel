package message

import (
	"errors"
	"io"
	"net"

	"github.com/hashicorp/yamux"
)

func ErrOther(err error) bool {
	return err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) && !errors.Is(err, yamux.ErrStreamClosed)
}
