package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/client"
	"github.com/starudream/secret-tunnel/constant"

	"github.com/starudream/secret-tunnel/internal/osx"
)

var rootCmd = &cobra.Command{
	Use:     "client",
	Version: constant.VERSION + " (" + constant.BIDTIME + ")",
	Run: func(cmd *cobra.Command, args []string) {
		app.Add(client.Start)
		err := app.OnceGo()
		if err != nil {
			log.Error().Msgf("client init error: %v", err)
		}
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	rootCmd.PersistentFlags().String("addr", "127.0.0.1:9797", "server address")
	osx.PE(config.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr")))

	rootCmd.PersistentFlags().String("key", "", "auth key")
	osx.PE(config.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key")))

	rootCmd.PersistentFlags().String("dns", "8.8.8.8", "dns server")
	osx.PE(config.BindPFlag("dns", rootCmd.PersistentFlags().Lookup("dns")))

	rootCmd.AddCommand(serviceCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
