package main

import (
	"github.com/spf13/cobra"

	"github.com/starudream/secret-tunnel/model"

	"github.com/starudream/secret-tunnel/internal/osx"
	"github.com/starudream/secret-tunnel/internal/tablew"
)

var (
	taskCmd = &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
		Args:  cobra.MinimumNArgs(1),
	}

	taskCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create task",
		PreRun: func(cmd *cobra.Command, args []string) {
			osx.PE(model.Init())
		},
		Run: func(cmd *cobra.Command, args []string) {
			client, err := model.GetClientById(taskClientId)
			osx.PE(err, tablew.PrintStruct(client))
			task, err := model.CreateTask(&model.Task{ClientId: taskClientId, Name: taskName, Addr: taskAddr})
			osx.PA(err, tablew.PrintStruct(task))
		},
	}

	taskListCmd = &cobra.Command{
		Use:   "list",
		Short: "List task",
		PreRun: func(cmd *cobra.Command, args []string) {
			osx.PE(model.Init())
		},
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := model.ListTaskByClientId(taskClientId)
			osx.PA(err, tablew.PrintStructs(tasks))
		},
	}

	taskClientId uint
	taskName     string
	taskAddr     string
)

func init() {
	taskCreateCmd.PersistentFlags().UintVar(&taskClientId, "client-id", 0, "which client the task belongs to")
	osx.PE(taskCreateCmd.MarkPersistentFlagRequired("client-id"))
	taskCreateCmd.PersistentFlags().StringVar(&taskName, "name", "", "task name")
	osx.PE(taskCreateCmd.MarkPersistentFlagRequired("name"))
	taskCreateCmd.PersistentFlags().StringVar(&taskAddr, "addr", "", "task address")
	osx.PE(taskCreateCmd.MarkPersistentFlagRequired("addr"))

	taskListCmd.PersistentFlags().UintVar(&taskClientId, "client-id", 0, "which client the task belongs to")

	taskCmd.AddCommand(taskCreateCmd)
	taskCmd.AddCommand(taskListCmd)
}
