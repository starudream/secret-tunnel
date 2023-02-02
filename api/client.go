package api

import (
	"strconv"

	"github.com/starudream/go-lib/errx"
	"github.com/starudream/go-lib/router"

	"github.com/starudream/secret-tunnel/model"
)

type clientReq struct {
	Name string `json:"name" validate:"required,max=24"`
	Key  string `json:"key,omitempty"`
}

func clientCreate(c *router.Context) {
	req := &clientReq{}
	if c.BindJSON(req) != nil {
		return
	}

	client, err := model.CreateClient(&model.Client{Name: req.Name})
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK(client)
}

func clientGet(c *router.Context) {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		return
	}

	client, err := model.GetClientById(uint(cid))
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK(client)
}

func clientList(c *router.Context) {
	clients, err := model.ListClient()
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK(clients)
}

func clientUpdate(c *router.Context) {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		c.Error(errx.ErrParam.WithMessage("invalid cid"))
		return
	}

	client, err := model.GetClientById(uint(cid))
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	req := &clientReq{}
	if c.BindJSON(req) != nil {
		return
	}

	client.Name = req.Name

	client, err = model.UpdateClient(client)
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK(client)
}

func clientDelete(c *router.Context) {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		c.Error(errx.ErrParam.WithMessage("invalid cid"))
		return
	}

	err = model.DeleteClient(uint(cid))
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	err = model.DeleteTaskByClientId(uint(cid))
	if err != nil {
		c.Error(errx.From(model.Wrap(err)))
		return
	}

	c.OK()
}
