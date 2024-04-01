package main

import (
	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/service/v2"
)

var rootCmd = cobra.NewRootCommand(func(c *cobra.Command) {
	c.Use = "server"

	c.PersistentFlags().String("addr", "0.0.0.0:9797", "server address")

	c.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		config.LoadFlags(c.PersistentFlags())
	}
	c.RunE = func(cmd *cobra.Command, args []string) error {
		return service.New("secret-tunnel-server", nil).Run()
	}

	cobra.AddConfigFlag(c)
})
