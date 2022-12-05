package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/starudream/secret-tunnel/message"
)

func messageSend(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	cid, err := strconv.Atoi(ps.ByName("cid"))
	if err != nil {
		ERRRequest(w, "invalid cid")
		return
	}

	switch strings.ToLower(ps.ByName("action")) {
	case "stop":
		COMM <- &message.StopService{Cid: uint(cid)}
	case "uninstall":
		COMM <- &message.UninstallService{Cid: uint(cid)}
	}

	OK(w)
}
