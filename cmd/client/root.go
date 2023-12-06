package main

import (
	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/config"

	"github.com/starudream/secret-tunnel/client"
)

var rootCmd = cobra.NewRootCommand(func(c *cobra.Command) {
	c.Use = "client"

	c.PersistentFlags().String("addr", "127.0.0.1:9797", "server address")
	c.PersistentFlags().String("dns", "119.29.29.29:53", "dns server")
	c.PersistentFlags().String("key", "", "auth key")

	c.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		config.LoadFlags(c.PersistentFlags())
	}
	c.RunE = func(cmd *cobra.Command, args []string) error {
		return runClient()
	}

	cobra.AddConfigFlag(c)
})

func runClient() error {
	return client.Run()
}
