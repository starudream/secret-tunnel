package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/starudream/secret-tunnel/message"
)

func messageSend(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	cid, err := strconv.Atoi(ps.ByName("cid"))
	if err != nil {
		ERRRequest(w, "invalid cid")
		return
	}

	COMM <- &message.UninstallService{Cid: uint(cid)}

	OK(w)
}
