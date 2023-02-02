package api

import (
	"strconv"

	"github.com/starudream/go-lib/errx"
	"github.com/starudream/go-lib/router"

	"github.com/starudream/secret-tunnel/model"
)

type taskReq struct {
	ClientId uint   `json:"client_id"`
	Name     string `json:"name" validate:"required,max=24"`
	Secret   string `json:"secret,omitempty"`
	Addr     string `json:"addr" validate:"required,hostname_port"`
}

func taskCreate(c *router.Context) {
	req := &taskReq{}
	if c.BindJSON(req) != nil {
		return
	}

	if req.ClientId <= 0 {
		c.Error(errx.ErrParam.WithMessage("missing client_id"))
		return
	}

	task, err := model.CreateTask(&model.Task{ClientId: req.ClientId, Name: req.Name, Addr: req.Addr})
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK(task)
}

func taskGet(c *router.Context) {
	tid, err := strconv.Atoi(c.Query("tid"))
	if err != nil {
		c.Error(errx.ErrParam.WithMessage("invalid tid"))
		return
	}

	task, err := model.GetTaskById(uint(tid))
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK(task)
}

func taskList(c *router.Context) {
	cid, _ := strconv.Atoi(c.Query("cid"))

	tasks, err := model.ListTaskByClientId(uint(cid))
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK(tasks)
}

func taskUpdate(c *router.Context) {
	tid, err := strconv.Atoi(c.Param("tid"))
	if err != nil {
		c.Error(errx.ErrParam.WithMessage("invalid tid"))
		return
	}

	task, err := model.GetTaskById(uint(tid))
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	req := &taskReq{}
	if c.BindJSON(req) != nil {
		return
	}

	task.Name = req.Name
	task.Addr = req.Addr

	task, err = model.UpdateTask(task)
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK(task)
}

func taskDelete(c *router.Context) {
	tid, err := strconv.Atoi(c.Query("tid"))
	if err != nil {
		c.Error(errx.ErrParam.WithMessage("invalid tid"))
		return
	}

	err = model.DeleteTask(uint(tid))
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK()
}
