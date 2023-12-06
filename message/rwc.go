package message

import (
	"io"
	"sync/atomic"

	"github.com/starudream/secret-tunnel/message/snappy"
)

type ReadWriteCloser struct {
	r io.Reader
	w io.Writer

	closed  atomic.Bool
	closeFn func() error
}

var _ io.ReadWriteCloser = (*ReadWriteCloser)(nil)

func (rwc *ReadWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.r.Read(p)
}

func (rwc *ReadWriteCloser) Write(p []byte) (n int, err error) {
	n, err = rwc.w.Write(p)
	if err != nil {
		return
	}
	if w, ok := rwc.w.(*snappy.Writer); ok {
		err = w.Flush()
	}
	return
}

func (rwc *ReadWriteCloser) Close() error {
	if !rwc.closed.CompareAndSwap(false, true) {
		return nil
	}

	if rc, ok := rwc.r.(io.Closer); ok {
		err := rc.Close()
		if err != nil {
			return err
		}
	}

	if wc, ok := rwc.w.(io.Closer); ok {
		err := wc.Close()
		if err != nil {
			return err
		}
	}

	if rwc.closeFn != nil {
		err := rwc.closeFn()
		if err != nil {
			return err
		}
	}

	return nil
}

func WrapReadWriteCloser(r io.Reader, w io.Writer, closeFn func() error) io.ReadWriteCloser {
	return &ReadWriteCloser{r: r, w: w, closeFn: closeFn}
}
