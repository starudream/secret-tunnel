package client

import (
	"context"

	"github.com/kardianos/service"

	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/log"
)

var Service = &iService{}

type iService struct {
}

var _ service.Interface = (*iService)(nil)

func (i *iService) Start(_ service.Service) error {
	go func() {
		err := Start(context.Background())
		if err != nil {
			log.Fatal().Msgf("client init error: %v", err)
		}
	}()
	return nil
}

func (i *iService) Stop(_ service.Service) error {
	app.Stop()
	return nil
}
