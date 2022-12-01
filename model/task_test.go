package model

import (
	"testing"

	"github.com/starudream/go-lib/randx"
	"github.com/starudream/go-lib/seq"
	"github.com/starudream/go-lib/testx"
)

func TestTask(t *testing.T) {
	client, err := CreateClient(&Client{Name: randx.F().Name(), Key: seq.UUID()})
	testx.P(t, err, client)

	task, err := CreateTask(&Task{
		ClientId: client.Id,
		Name:     randx.F().Name(),
		Secret:   seq.UUID(),
		Addr:     randx.F().IPv4Address(),
	})
	testx.P(t, err, task)

	task, err = GetTaskById(task.Id)
	testx.P(t, err, task)

	task.Name = randx.F().Name()

	_, err = UpdateTask(task)
	testx.P(t, err)

	err = UpdateTaskActive(task.Id, true)
	testx.P(t, err)

	task, err = GetTaskBySecret(client.Id, task.Secret)
	testx.P(t, err, task)

	tasks, err := ListTaskByClientId(client.Id)
	testx.P(t, err, tasks)

	err = DeleteTask(task.Id)
	testx.P(t, err)
}
