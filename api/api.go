package api

import (
	"context"
	"net"
	"net/http"

	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/errx"
	"github.com/starudream/go-lib/log"
	"github.com/starudream/go-lib/router"

	"github.com/starudream/secret-tunnel/constant"
)

type M map[string]any

var COMM chan any

var server *http.Server

func Start(ctx context.Context) error {
	register()

	address := config.GetString("api")

	server = &http.Server{Addr: address, Handler: router.Handler()}

	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Info().Msgf("api start success on %s", address)

	return server.Serve(ln)
}

func Stop() {
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Error().Msgf("api shutdown error: %v", err)
	}

	log.Info().Msgf("api shutdown")
}

func register() {
	token := config.GetString("token")
	if token == "" {
		log.Warn().Msgf("empty token cause security risks")
	} else {
		router.Use(auth(token))
	}

	router.Handle(http.MethodGet, "/", index)
	router.Handle(http.MethodPost, "/", index)

	router.Handle(http.MethodPost, "/client", clientCreate)
	router.Handle(http.MethodGet, "/client/:cid", clientGet)
	router.Handle(http.MethodGet, "/clients", clientList)
	router.Handle(http.MethodPatch, "/client/:cid", clientUpdate)
	router.Handle(http.MethodDelete, "/client/:cid", clientDelete)

	router.Handle(http.MethodPost, "/task", taskCreate)
	router.Handle(http.MethodGet, "/task/:tid", taskGet)
	router.Handle(http.MethodGet, "/tasks", taskList)
	router.Handle(http.MethodPatch, "/task/:tid", taskUpdate)
	router.Handle(http.MethodDelete, "/task/:tid", taskDelete)

	router.Handle(http.MethodPost, "/message/:cid/:action", messageSend)
}

func auth(token string) func(*router.Context) {
	return func(c *router.Context) {
		s := stringNotEmpty(c.GetHeader("Authorization"), c.GetHeader("Token"), c.Query("token"))
		if s == "" {
			c.AbortWithError(errx.ErrUnAuth.WithMessage("missing token"))
			return
		}
		if s != token {
			c.AbortWithError(errx.ErrUnAuth.WithMessage("token not match"))
			return
		}
		c.Next()
	}
}

func index(c *router.Context) {
	c.OK(M{"version": constant.VERSION, "bidtime": constant.BIDTIME})
}

func stringNotEmpty(vs ...string) string {
	for i := 0; i < len(vs); i++ {
		if vs[i] != "" {
			return vs[i]
		}
	}
	return ""
}
