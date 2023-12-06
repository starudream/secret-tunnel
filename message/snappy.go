package message

import (
	"io"

	"github.com/starudream/secret-tunnel/message/snappy"
)

var snappyPool = snappy.Pool{}

func WithSnappy(rwc io.ReadWriteCloser) io.ReadWriteCloser {
	sr, sw := snappyPool.GetReader(rwc), snappyPool.GetWriter(rwc)
	return WrapReadWriteCloser(sr, sw, func() error {
		err := rwc.Close()
		snappyPool.PutReader(sr)
		snappyPool.PutWriter(sw)
		return err
	})
}
