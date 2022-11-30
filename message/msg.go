package message

import (
	"fmt"
	"reflect"
)

// a b c d e f g
// h i j k l m n
// o p q r s t
// u v w x y z

var (
	typeMap = map[byte]any{
		// A B means common
		0xA1: &OK{},
		0xA2: &Data{},
		0xA3: &Close{},
		0xA4: &Error{},
		0xA5: &Task{},

		// C D means server -> client
		0xC1: &LoginResp{},
		0xC2: &WorkResp{},
		0xC3: &CreateTaskReq{},
		0xC4: &ConnectTaskResp{},
		0xC5: &CloseTaskReq{},

		// E F means server -> client
		0xE1: &LoginReq{},
		0xE2: &WorkReq{},
		0xE3: &CreateTaskResp{},
		0xE4: &ConnectTaskReq{},
		0xE5: &CloseTaskResp{},
	}

	typeByteMap = map[reflect.Type]byte{}
	byteTypeMap = map[byte]reflect.Type{}
)

func init() {
	for b, a := range typeMap {
		t := reflect.TypeOf(a)
		typeByteMap[t] = b
		byteTypeMap[b] = t.Elem()
	}
}

// common

type OK struct {
}

type Data struct {
}

type Close struct {
	Reason string `json:"reason,omitempty"`
}

type Error struct {
	Error string `json:"error,omitempty"`
}

func NewError(v string, a ...any) Error {
	if len(a) == 0 {
		return Error{Error: v}
	}
	return Error{Error: fmt.Sprintf(v, a...)}
}

func (err Error) GetError() string {
	return err.Error
}

type Task struct {
	Id     uint   `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Secret string `json:"secret,omitempty"`
}

func NewTask(id uint, name, addr string) Task {
	return Task{
		Id:   id,
		Name: name,
		Addr: addr,
	}
}

func NewTaskSecret(secret string) Task {
	return Task{
		Secret: secret,
	}
}

// server -> client

type LoginResp struct {
	Ver string `json:"ver,omitempty"`

	Error
}

type WorkResp struct {
	Wid string `json:"wid,omitempty"`

	Error
}

type CreateTaskReq struct {
	Tid string `json:"tid,omitempty"`
	Sid string `json:"sid,omitempty"`

	Task
}

type ConnectTaskResp struct {
	Sid string `json:"sid,omitempty"`

	Task
	Error
}

type CloseTaskReq struct {
	Tid string `json:"tid,omitempty"`
}

// client -> server

type LoginReq struct {
	Ver string `json:"ver,omitempty"`
	Key string `json:"key,omitempty"`

	GO       string `json:"go,omitempty"`
	OS       string `json:"os,omitempty"`
	ARCH     string `json:"arch,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

type WorkReq struct {
}

type CreateTaskResp struct {
	Sid string `json:"sid,omitempty"`

	Error
}

type ConnectTaskReq struct {
	Wid string `json:"wid,omitempty"`
	Tid string `json:"tid,omitempty"`

	Task
}

type CloseTaskResp struct {
}
