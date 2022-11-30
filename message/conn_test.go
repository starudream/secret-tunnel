package message

import (
	"net"
	"runtime"
	"testing"

	"github.com/starudream/go-lib/testx"
)

func TestMessage(t *testing.T) {
	server, client := net.Pipe()

	defer func() {
		_ = server.Close()
		_ = client.Close()
	}()

	done := make(chan struct{})

	go func() {
		testx.RequireNoErrorf(t, Write(server, LoginReq{GO: runtime.Version(), OS: runtime.GOOS, ARCH: runtime.GOARCH}), "1")
		testx.RequireNoErrorf(t, Write(server, Close{}), "close")
	}()

	go func() {
		for {
			v, ve := Read(client)
			testx.P(t, ve, v)

			_, ok := v.(*Close)
			if ok {
				close(done)
				break
			}
		}
	}()

	<-done
}
