package message

import (
	"net"
	"runtime"
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestMessage(t *testing.T) {
	server, client := net.Pipe()

	defer func() {
		_ = server.Close()
		_ = client.Close()
	}()

	done := make(chan struct{})

	go func() {
		testutil.LogNoErr(t, Write(server, LoginReq{GO: runtime.Version(), OS: runtime.GOOS, ARCH: runtime.GOARCH}), "1")
		testutil.LogNoErr(t, Write(server, Close{}), "close")
	}()

	go func() {
		for {
			v, ve := Read(client)
			testutil.LogNoErr(t, ve, v)

			_, ok := v.(*Close)
			if ok {
				close(done)
				break
			}
		}
	}()

	<-done
}
