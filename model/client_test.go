package model

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/starudream/go-lib/core/v2/config/version"
	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestClient(t *testing.T) {
	client, err := CreateClient(&Client{Name: "CA" + time.Now().Format("20060102150405")})
	testutil.LogNoErr(t, err, client)

	client, err = GetClientById(client.Id)
	testutil.LogNoErr(t, err, client)

	client.Name = "CB" + time.Now().Format("20060102150405")

	_, err = UpdateClient(client)
	testutil.LogNoErr(t, err)

	err = UpdateClientActive(client.Id, false)
	testutil.LogNoErr(t, err)

	err = UpdateClientOnline(&Client{
		Id:       client.Id,
		Ver:      version.GetVersionInfo().GitVersion,
		Addr:     "127.0.0.1",
		GO:       runtime.Version(),
		OS:       runtime.GOOS,
		ARCH:     runtime.GOARCH,
		Hostname: func() string { name, _ := os.Hostname(); return name }(),
	})
	testutil.LogNoErr(t, err)

	client, err = GetClientByKey(client.Key)
	testutil.LogNoErr(t, err, client)

	err = UpdateClientOffline(client.Id)
	testutil.LogNoErr(t, err)

	clients, err := ListClient()
	testutil.LogNoErr(t, err, clients)

	err = DeleteClient(client.Id)
	testutil.LogNoErr(t, err)
}
