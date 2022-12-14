package pool

import (
	"io"
	"sync"

	"github.com/golang/snappy"
)

var (
	snappyReaderPool sync.Pool
	snappyWriterPool sync.Pool
)

func GetSnappyReader(r io.Reader) *snappy.Reader {
	x := snappyReaderPool.Get()
	if x == nil {
		return snappy.NewReader(r)
	}
	sr := x.(*snappy.Reader)
	sr.Reset(r)
	return sr
}

func PutSnappyReader(sr *snappy.Reader) {
	snappyReaderPool.Put(sr)
}

func GetSnappyWriter(w io.Writer) *snappy.Writer {
	x := snappyWriterPool.Get()
	if x == nil {
		return snappy.NewBufferedWriter(w)
	}
	sw := x.(*snappy.Writer)
	sw.Reset(w)
	return sw
}

func PutSnappyWriter(sw *snappy.Writer) {
	snappyWriterPool.Put(sw)
}
