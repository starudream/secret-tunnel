package main

import (
	"github.com/spf13/cobra"

	"github.com/starudream/secret-tunnel/model"

	"github.com/starudream/secret-tunnel/internal/osx"
	"github.com/starudream/secret-tunnel/internal/tablew"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Manage clients",
		Args:  cobra.MinimumNArgs(1),
	}

	clientCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create client",
		PreRun: func(cmd *cobra.Command, args []string) {
			osx.PE(model.Init())
		},
		Run: func(cmd *cobra.Command, args []string) {
			client, err := model.CreateClient(&model.Client{Name: clientName})
			osx.PA(err, tablew.PrintStruct(client))
		},
	}

	clientListCmd = &cobra.Command{
		Use:   "list",
		Short: "List client",
		PreRun: func(cmd *cobra.Command, args []string) {
			osx.PE(model.Init())
		},
		Run: func(cmd *cobra.Command, args []string) {
			clients, err := model.ListClient()
			osx.PA(err, tablew.PrintStructs(clients))
		},
	}

	clientName string
)

func init() {
	clientCreateCmd.PersistentFlags().StringVar(&clientName, "name", "", "client name")
	osx.PE(clientCreateCmd.MarkPersistentFlagRequired("name"))

	clientCmd.AddCommand(clientCreateCmd)
	clientCmd.AddCommand(clientListCmd)
}
