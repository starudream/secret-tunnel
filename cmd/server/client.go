package main

import (
	"fmt"

	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/utils/sliceutil"

	"github.com/starudream/secret-tunnel/model"
	"github.com/starudream/secret-tunnel/util"
)

var (
	clientCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "client"
		c.Short = "Manage client"
	})

	clientListCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "list"
		c.Aliases = []string{"ls"}
		c.Short = "List client"
		c.RunE = func(cmd *cobra.Command, args []string) error {
			clients, err := model.ListClient()
			if err != nil {
				return err
			}
			return util.TablePrint(clients)
		}
	})

	clientCreateCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "create <name>"
		c.Short = "Create client"
		c.RunE = func(cmd *cobra.Command, args []string) error {
			name, _ := sliceutil.GetValue(args, 0)
			if name == "" {
				return fmt.Errorf("requires client name")
			}
			client, err := model.CreateClient(&model.Client{Name: name})
			if err != nil {
				return err
			}
			return util.TablePrint([]*model.Client{client})
		}
	})
)

func init() {
	clientCmd.AddCommand(clientListCmd)
	clientCmd.AddCommand(clientCreateCmd)

	rootCmd.AddCommand(clientCmd)
}
