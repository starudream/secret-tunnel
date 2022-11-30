package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/api"
	"github.com/starudream/secret-tunnel/constant"
	"github.com/starudream/secret-tunnel/model"
	"github.com/starudream/secret-tunnel/server"
)

var rootCmd = &cobra.Command{
	Use:     constant.AppName,
	Version: constant.VERSION + " (" + constant.BIDTIME + ")",
	RunE: func(cmd *cobra.Command, args []string) error {
		app.Init(model.Init)
		app.Add(api.Start, server.Start)
		err := app.OnceGo()
		if err != nil {
			log.Error().Msgf("server init error: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "version of "+constant.AppName)

	rootCmd.PersistentFlags().String("addr", "0.0.0.0:9797", "server address")
	_ = config.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))

	rootCmd.PersistentFlags().String("api", "127.0.0.1:9799", "api address")
	_ = config.BindPFlag("api", rootCmd.PersistentFlags().Lookup("api"))

	rootCmd.PersistentFlags().String("token", "", "api token")
	_ = config.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
