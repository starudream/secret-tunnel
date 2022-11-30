package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/starudream/go-lib/seq"

	"github.com/starudream/secret-tunnel/model"
)

type clientReq struct {
	Name string `json:"name" validate:"required,max=24"`
	Key  string `json:"key,omitempty"`
}

func clientCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req, err := V[*clientReq](w, r)
	if err != nil {
		return
	}

	req.Key = seq.UUIDShort()

	client, err := model.CreateClient(&model.Client{Name: req.Name, Key: req.Key})
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w, client)
}

func clientGet(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	cid, err := strconv.Atoi(ps.ByName("cid"))
	if err != nil {
		ERRRequest(w, "invalid cid")
		return
	}

	client, err := model.GetClientById(uint(cid))
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w, client)
}

func clientList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	clients, err := model.ListClient()
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w, clients)
}

func clientUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cid, err := strconv.Atoi(ps.ByName("cid"))
	if err != nil {
		ERRRequest(w, "invalid cid")
		return
	}

	client, err := model.GetClientById(uint(cid))
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	req, err := V[*clientReq](w, r)
	if err != nil {
		return
	}

	client.Name = req.Name

	client, err = model.UpdateClient(client)
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w, client)
}

func clientDelete(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	cid, err := strconv.Atoi(ps.ByName("cid"))
	if err != nil {
		ERRRequest(w, "invalid cid")
		return
	}

	err = model.DeleteClient(uint(cid))
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w)
}
