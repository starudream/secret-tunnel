package netx

import (
	"io"
	"sync"
)

type ReadWriteCloser struct {
	r io.Reader
	w io.Writer

	closed  bool
	closeFn func() error

	mu sync.Mutex
}

var _ io.ReadWriteCloser = (*ReadWriteCloser)(nil)

func (rwc *ReadWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.r.Read(p)
}

func (rwc *ReadWriteCloser) Write(p []byte) (n int, err error) {
	return rwc.w.Write(p)
}

func (rwc *ReadWriteCloser) Close() (res error) {
	rwc.mu.Lock()
	if rwc.closed {
		rwc.mu.Unlock()
		return
	}
	rwc.closed = true
	rwc.mu.Unlock()

	var err error
	if rc, ok := rwc.r.(io.Closer); ok {
		err = rc.Close()
		if err != nil {
			res = err
		}
	}

	if wc, ok := rwc.w.(io.Closer); ok {
		err = wc.Close()
		if err != nil {
			res = err
		}
	}

	if rwc.closeFn != nil {
		err = rwc.closeFn()
		if err != nil {
			res = err
		}
	}
	return
}

func WrapReadWriteCloser(r io.Reader, w io.Writer, closeFn func() error) io.ReadWriteCloser {
	return &ReadWriteCloser{
		r:       r,
		w:       w,
		closeFn: closeFn,
	}
}
