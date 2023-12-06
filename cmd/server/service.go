package main

import (
	"context"

	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/service/v2"
)

func init() {
	service.AddCommand(rootCmd, service.New("secret-tunnel-server", serviceRun, service.WithArguments()))
}

func serviceRun(context.Context) {
	err := runServer()
	if err != nil {
		slog.Error("server run error: %v", err)
	}
}
