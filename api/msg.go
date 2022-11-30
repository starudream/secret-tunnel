package api

import (
	"fmt"
	"net/http"

	"github.com/starudream/go-lib/codec/json"
)

type Resp struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
}

func NewResp(code int, s string, v ...any) *Resp {
	resp := &Resp{Code: code}
	if len(v) == 0 {
		resp.Error = s
	} else {
		resp.Error = fmt.Sprintf(s, v...)
	}
	return resp
}

type M map[string]any

func OK(w http.ResponseWriter, v ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(v) > 0 && v[0] != nil {
		_, _ = w.Write(json.MustMarshal(v[0]))
	} else {
		_, _ = w.Write([]byte(`{"msg":"ok"}`))
	}
}

func ERR(w http.ResponseWriter, code int, s string, v ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(json.MustMarshal(NewResp(code, s, v...)))
}

func ERRInternal(w http.ResponseWriter, s string, v ...any) {
	ERR(w, http.StatusInternalServerError, s, v...)
}

func ERRRequest(w http.ResponseWriter, s string, v ...any) {
	ERR(w, http.StatusBadRequest, s, v...)
}
