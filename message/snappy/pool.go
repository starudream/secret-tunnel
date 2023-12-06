package snappy

import (
	"io"
	"sync"

	"github.com/golang/snappy"
)

type (
	Reader = snappy.Reader
	Writer = snappy.Writer
)

type Pool struct {
	rs sync.Pool
	ws sync.Pool
}

func (p *Pool) GetReader(r io.Reader) io.Reader {
	if v := p.rs.Get(); v != nil {
		sr := v.(*Reader)
		sr.Reset(r)
		return sr
	}
	return snappy.NewReader(r)
}

func (p *Pool) PutReader(r io.Reader) {
	p.rs.Put(r)
}

func (p *Pool) GetWriter(w io.Writer) io.Writer {
	if v := p.ws.Get(); v != nil {
		sw := v.(*Writer)
		sw.Reset(w)
		return sw
	}
	return snappy.NewBufferedWriter(w)
}

func (p *Pool) PutWriter(w io.Writer) {
	p.ws.Put(w)
}
