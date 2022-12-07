package model

import (
	"os"
	"runtime"
	"testing"

	"github.com/starudream/go-lib/randx"
	"github.com/starudream/go-lib/testx"
)

func TestClient(t *testing.T) {
	client, err := CreateClient(&Client{Name: randx.F().Name()})
	testx.P(t, err, client)

	client, err = GetClientById(client.Id)
	testx.P(t, err, client)

	client.Name = randx.F().Name()

	_, err = UpdateClient(client)
	testx.P(t, err)

	err = UpdateClientActive(client.Id, true)
	testx.P(t, err)

	err = UpdateClientOnline(&Client{
		Id:       client.Id,
		Addr:     randx.F().IPv4Address(),
		GO:       runtime.Version(),
		OS:       runtime.GOOS,
		ARCH:     runtime.GOARCH,
		Hostname: func() string { name, _ := os.Hostname(); return name }(),
	})
	testx.P(t, err)

	client, err = GetClientByKey(client.Key)
	testx.P(t, err, client)

	err = UpdateClientOffline(client.Id)
	testx.P(t, err)

	clients, err := ListClient()
	testx.P(t, err, clients)

	err = DeleteClient(client.Id)
	testx.P(t, err)
}
