package message

import (
	"io"
	"sync"

	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/poolutil"
)

func Copy(c1, c2 io.ReadWriteCloser) (in, out int64) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		in = copyFn(c1, c2)
	}()
	go func() {
		defer wg.Done()
		out = copyFn(c2, c1)
	}()
	wg.Wait()
	return
}

var copyPool = poolutil.NewBytes(100, 1024*MaxSize)

func copyFn(dst, src io.ReadWriteCloser) int64 {
	buf := copyPool.Get()
	defer copyPool.Put(buf)

	defer gh.Close(src)
	defer gh.Close(dst)

	total, err := io.CopyBuffer(dst, src, buf)
	if ErrOther(err) {
		slog.Warn("copy buffer error: %v", err)
	}

	return total
}
