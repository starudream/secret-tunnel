package model

import (
	"testing"
	"time"

	"github.com/starudream/go-lib/core/v2/utils/testutil"
)

func TestTask(t *testing.T) {
	client, err := CreateClient(&Client{Name: "CC" + time.Now().Format("20060102150405")})
	testutil.LogNoErr(t, err, client)

	task, err := CreateTask(&Task{
		ClientId: client.Id,
		Name:     "CA" + time.Now().Format("20060102150405"),
		Addr:     "127.0.0.1",
	})
	testutil.LogNoErr(t, err, task)

	task, err = GetTaskById(task.Id)
	testutil.LogNoErr(t, err, task)

	task.Name = "CB" + time.Now().Format("20060102150405")

	_, err = UpdateTask(task)
	testutil.LogNoErr(t, err)

	err = UpdateTaskActive(task.Id, false)
	testutil.LogNoErr(t, err)

	err = UpdateTaskCompress(task.Id, true)
	testutil.LogNoErr(t, err)

	err = UpdateTaskTraffic(task.Id, 100, 200)
	testutil.LogNoErr(t, err)

	task, err = GetTaskBySecret(client.Id, task.Secret)
	testutil.LogNoErr(t, err, task)

	tasks, err := ListTaskByClientId(client.Id)
	testutil.LogNoErr(t, err, tasks)

	err = DeleteTask(task.Id)
	testutil.LogNoErr(t, err)

	err = DeleteTaskByClientId(client.Id)
	testutil.LogNoErr(t, err)
}
