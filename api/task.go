package api

import (
	"net/http"
	"strconv"

	"github.com/starudream/secret-tunnel/model"
)

type taskReq struct {
	ClientId uint   `json:"client_id"`
	Name     string `json:"name" validate:"required,max=24"`
	Secret   string `json:"secret,omitempty"`
	Addr     string `json:"addr" validate:"required,hostname_port"`
}

func taskCreate(w http.ResponseWriter, r *http.Request, _ Params) {
	req, err := V[*taskReq](w, r)
	if err != nil {
		return
	}

	if req.ClientId <= 0 {
		ERRRequest(w, "missing client_id")
		return
	}

	task, err := model.CreateTask(&model.Task{ClientId: req.ClientId, Name: req.Name, Addr: req.Addr})
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w, task)
}

func taskGet(w http.ResponseWriter, _ *http.Request, ps Params) {
	tid, err := strconv.Atoi(ps.ByName("tid"))
	if err != nil {
		ERRRequest(w, "invalid tid")
		return
	}

	task, err := model.GetTaskById(uint(tid))
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w, task)
}

func taskList(w http.ResponseWriter, r *http.Request, _ Params) {
	cid, _ := strconv.Atoi(r.URL.Query().Get("cid"))

	tasks, err := model.ListTaskByClientId(uint(cid))
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w, tasks)
}

func taskUpdate(w http.ResponseWriter, r *http.Request, ps Params) {
	tid, err := strconv.Atoi(ps.ByName("tid"))
	if err != nil {
		ERRRequest(w, "invalid tid")
		return
	}

	task, err := model.GetTaskById(uint(tid))
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	req, err := V[*taskReq](w, r)
	if err != nil {
		return
	}

	task.Name = req.Name
	task.Addr = req.Addr

	task, err = model.UpdateTask(task)
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w, task)
}

func taskDelete(w http.ResponseWriter, _ *http.Request, ps Params) {
	tid, err := strconv.Atoi(ps.ByName("tid"))
	if err != nil {
		ERRRequest(w, "invalid tid")
		return
	}

	err = model.DeleteTask(uint(tid))
	if err != nil {
		ERRInternal(w, model.Wrap(err).Error())
		return
	}

	OK(w)
}
