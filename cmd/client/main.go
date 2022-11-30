package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/client"
	"github.com/starudream/secret-tunnel/constant"
)

var rootCmd = &cobra.Command{
	Use:     constant.AppName,
	Version: constant.VERSION + " (" + constant.BIDTIME + ")",
	RunE: func(cmd *cobra.Command, args []string) error {
		app.Add(client.Start)
		err := app.OnceGo()
		if err != nil {
			log.Error().Msgf("client init error: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "version of "+constant.AppName)

	rootCmd.PersistentFlags().String("addr", "", "server address")
	_ = config.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))

	rootCmd.PersistentFlags().String("key", "", "auth key")
	_ = config.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))

	rootCmd.PersistentFlags().String("dns", "8.8.8.8", "dns server")
	_ = config.BindPFlag("dns", rootCmd.PersistentFlags().Lookup("dns"))

	rootCmd.PersistentFlags().String("tasks", "[]", "tasks json string")
	_ = config.BindPFlag("tasks", rootCmd.PersistentFlags().Lookup("tasks"))
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
