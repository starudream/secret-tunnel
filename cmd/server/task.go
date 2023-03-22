package main

import (
	"github.com/starudream/go-lib/flag"

	"github.com/starudream/secret-tunnel/model"

	"github.com/starudream/secret-tunnel/internal/osx"
	"github.com/starudream/secret-tunnel/internal/tablew"
)

var (
	taskCmd = &flag.Command{
		Use:   "task",
		Short: "Manage tasks",
		Args:  flag.MinimumNArgs(1),
	}

	taskCreateCmd = &flag.Command{
		Use:   "create",
		Short: "Create task",
		PreRun: func(cmd *flag.Command, args []string) {
			osx.PE(model.Init())
		},
		Run: func(cmd *flag.Command, args []string) {
			client, err := model.GetClientById(taskClientId)
			osx.PE(err, tablew.PrintStruct(client))
			task, err := model.CreateTask(&model.Task{ClientId: taskClientId, Name: taskName, Addr: taskAddr})
			osx.PA(err, tablew.PrintStruct(task))
		},
	}

	taskListCmd = &flag.Command{
		Use:   "list",
		Short: "List task",
		PreRun: func(cmd *flag.Command, args []string) {
			osx.PE(model.Init())
		},
		Run: func(cmd *flag.Command, args []string) {
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
