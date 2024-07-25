package main

import (
	"context"

	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/service/v2"

	"github.com/starudream/secret-tunnel/client"
)

func init() {
	args := cobra.FlagArgs(rootCmd.PersistentFlags(), "config")
	if c := config.LoadedFile(); c != "" {
		args = append(args, "-c", c)
	}
	service.AddCommand(rootCmd, service.New("secret-tunnel-client", serviceRun, service.WithArguments(args...)))
}

func serviceRun(context.Context) {
	err := client.Run()
	if err != nil {
		slog.Fatal("client run error: %v", err)
	}
}
