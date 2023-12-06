package main

import (
	"fmt"
	"strconv"

	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/utils/sliceutil"

	"github.com/starudream/secret-tunnel/model"
	"github.com/starudream/secret-tunnel/util"
)

var (
	taskCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "task"
		c.Short = "Manage task"
	})

	taskListCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "list"
		c.Aliases = []string{"ls"}
		c.Short = "List task"
		c.RunE = func(cmd *cobra.Command, args []string) error {
			clientIdStr, _ := sliceutil.GetValue(args, 0)
			clientId, _ := strconv.Atoi(clientIdStr)
			tasks, err := model.ListTaskByClientId(uint(clientId))
			if err != nil {
				return err
			}
			return util.TablePrint(tasks)
		}
	})

	taskCreateCmd = cobra.NewCommand(func(c *cobra.Command) {
		c.Use = "create <client id> <name> <addr>"
		c.Short = "Create task"
		c.RunE = func(cmd *cobra.Command, args []string) error {
			clientIdStr, _ := sliceutil.GetValue(args, 0)
			if clientIdStr == "" {
				return fmt.Errorf("requires client id")
			}
			clientId, err := strconv.Atoi(clientIdStr)
			if err != nil {
				return err
			}
			name, _ := sliceutil.GetValue(args, 1)
			if name == "" {
				return fmt.Errorf("requires task name")
			}
			addr, _ := sliceutil.GetValue(args, 2)
			if addr == "" {
				return fmt.Errorf("requires task address")
			}
			task, err := model.CreateTask(&model.Task{ClientId: uint(clientId), Name: name, Addr: addr})
			if err != nil {
				return err
			}
			return util.TablePrint([]*model.Task{task})
		}
	})
)

func init() {
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskCreateCmd)

	rootCmd.AddCommand(taskCmd)
}
