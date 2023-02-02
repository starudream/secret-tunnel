package api

import (
	"strconv"
	"strings"

	"github.com/starudream/go-lib/errx"
	"github.com/starudream/go-lib/router"

	"github.com/starudream/secret-tunnel/message"
)

func messageSend(c *router.Context) {
	cid, err := strconv.Atoi(c.Query("cid"))
	if err != nil {
		c.Error(errx.ErrParam.WithMessage("invalid cid"))
		return
	}

	switch strings.ToLower(c.Param("action")) {
	case "stop":
		COMM <- &message.StopService{Cid: uint(cid)}
	case "uninstall":
		COMM <- &message.UninstallService{Cid: uint(cid)}
	}

	c.OK()
}
