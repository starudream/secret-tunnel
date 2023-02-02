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

	"github.com/starudream/secret-tunnel/internal/osx"
)

var rootCmd = &cobra.Command{
	Use:     "server",
	Version: constant.VERSION + " (" + constant.BIDTIME + ")",
	Run: func(cmd *cobra.Command, args []string) {
		app.Init(model.Init)
		app.Add(api.Start, server.Start)
		app.Defer(api.Stop)
		comm := make(chan any, 100)
		api.COMM, server.COMM = comm, comm
		err := app.OnceGo()
		if err != nil {
			log.Error().Msgf("server init error: %v", err)
		}
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	rootCmd.PersistentFlags().String("addr", "0.0.0.0:9797", "server address")
	osx.PE(config.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr")))

	rootCmd.PersistentFlags().String("api", "127.0.0.1:9799", "api address")
	osx.PE(config.BindPFlag("api", rootCmd.PersistentFlags().Lookup("api")))

	rootCmd.PersistentFlags().String("token", "", "api token")
	osx.PE(config.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token")))

	rootCmd.AddCommand(clientCmd)
	rootCmd.AddCommand(taskCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
