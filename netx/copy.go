package netx

import (
	"io"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"

	"github.com/starudream/go-lib/errx"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/constant"
	"github.com/starudream/secret-tunnel/pool"
)

func Copy(c1, c2 io.ReadWriteCloser) (in int64, out int64) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go connCopy(c1, c2, &in, wg)
	go connCopy(c2, c1, &out, wg)
	wg.Wait()
	log.Debug().Msgf("copy stats, in: %d, out: %d", in, out)
	return
}

func connCopy(dst, src io.ReadWriteCloser, n *int64, wg *sync.WaitGroup) {
	_ = copyPool.Invoke(&pfCopyItem{dst: dst, src: src, n: n, wg: wg})
}

var copyPool, _ = ants.NewPoolWithFunc(10000, pfCopy, ants.WithNonblocking(false))

type pfCopyItem struct {
	dst io.ReadWriteCloser
	src io.ReadWriteCloser

	n *int64

	wg *sync.WaitGroup
}

func (c *pfCopyItem) Done() {
	_ = c.dst.Close()
	_ = c.src.Close()
	c.wg.Done()
}

func pfCopy(v any) {
	i, ok := v.(*pfCopyItem)
	if !ok {
		return
	}

	defer i.Done()

	buf := pool.GetBuf(constant.MessageSize)
	defer pool.PutBuf(buf)

	var err error
	*i.n, err = io.CopyBuffer(i.dst, i.src, buf)
	if err != nil {
		if !errx.Is(err, io.EOF) && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Warn().Msgf("copy buffer error: %v", err)
		}
	}
}
