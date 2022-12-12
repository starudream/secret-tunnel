package api

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/starudream/go-lib/codec/json"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/constant"

	"github.com/starudream/secret-tunnel/internal/httpw"
)

var COMM chan any

func Start(ctx context.Context) error {
	a, err := newAPI(ctx)
	if err != nil {
		return err
	}

	a.register()

	s := &http.Server{Addr: a.address, Handler: a.router}

	ln, err := net.Listen("tcp", a.address)
	if err != nil {
		return err
	}

	log.Info().Msgf("api start success on %s", a.address)

	return s.Serve(ln)
}

type API struct {
	address string
	token   string

	ctx    context.Context
	cancel context.CancelFunc

	router *httprouter.Router
}

func newAPI(ctx context.Context) (*API, error) {
	ctx, cancel := context.WithCancel(ctx)
	a := &API{
		address: config.GetString("api"),
		token:   config.GetString("token"),
		ctx:     ctx,
		cancel:  cancel,
		router:  httprouter.New(),
	}
	if a.token == "" {
		log.Warn().Msgf("empty token cause security risks")
	}
	return a, nil
}

func (a *API) register() {
	a.router.HandleOPTIONS = true
	a.router.HandleMethodNotAllowed = true
	a.router.GlobalOPTIONS = handleOPTIONS()
	a.router.MethodNotAllowed = handleNotAllowed()
	a.router.NotFound = handleNotFound()
	a.router.PanicHandler = handlePanic()

	a.handle(http.MethodGet, "/", index)
	a.handle(http.MethodPost, "/", index)

	a.handleM(http.MethodPost, "/client", clientCreate)
	a.handleM(http.MethodGet, "/client/:cid", clientGet)
	a.handleM(http.MethodGet, "/clients", clientList)
	a.handleM(http.MethodPatch, "/client/:cid", clientUpdate)
	a.handleM(http.MethodDelete, "/client/:cid", clientDelete)

	a.handleM(http.MethodPost, "/task", taskCreate)
	a.handleM(http.MethodGet, "/task/:tid", taskGet)
	a.handleM(http.MethodGet, "/tasks", taskList)
	a.handleM(http.MethodPatch, "/task/:tid", taskUpdate)
	a.handleM(http.MethodDelete, "/task/:tid", taskDelete)

	a.handleM(http.MethodPost, "/message/:cid/:action", messageSend)
}

type Params = httprouter.Params

type middleware func(handle httprouter.Handle) httprouter.Handle

func (a *API) handle(method, path string, handle httprouter.Handle, middlewares ...middleware) {
	if ml := len(middlewares); ml > 0 {
		for i := ml - 1; i >= 0; i-- {
			handle = middlewares[i](handle)
		}
	}
	a.router.Handle(method, path, handle)
}

func (a *API) handleM(method, path string, handle httprouter.Handle) {
	ms := []middleware{a.logger}
	if a.token != "" {
		ms = append(ms, a.auth)
	}
	a.handle(method, path, handle, ms...)
}

func (a *API) auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		s := stringNotEmpty(r.Header.Get("Authorization"), r.Header.Get("Token"), r.URL.Query().Get("token"))
		if s == "" {
			ERR(w, http.StatusUnauthorized, "missing token")
			return
		}
		if s != a.token {
			ERR(w, http.StatusUnauthorized, "token not match")
			return
		}
		h(w, r, ps)
	}
}

func (a *API) logger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().Msgf("read http body error: %v", err)
			ERR(w, http.StatusInternalServerError, "internal server error")
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(bs))
		l := log.With().Str("remote", r.RemoteAddr).Str("method", r.Method).Str("path", r.URL.String()).Logger()
		var req any
		if json.Unmarshal(bs, &req) == nil {
			l.Info().Msgf("req=%s", json.MustMarshal(req))
		} else {
			l.Info().Msgf("req=%s", bs)
		}
		w = httpw.NewResponse(w)
		h(w, r, ps)
		resp, sc := httpw.GetResponse(w)
		l.Info().Int("code", sc).Dur("took", time.Since(start)).Msgf("resp=%s", resp)
	}
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	OK(w, M{"version": constant.VERSION, "bidtime": constant.BIDTIME})
}

func handleOPTIONS() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATH, DELETE, HEAD, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "43200")
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func handleNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ERR(w, http.StatusMethodNotAllowed, "method not allowed")
	})
}

func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ERR(w, http.StatusNotFound, "not found")
	})
}

func handlePanic() func(w http.ResponseWriter, r *http.Request, rcv any) {
	return func(w http.ResponseWriter, r *http.Request, rcv any) {
		log.Error().Msgf("api server panic: %s", debug.Stack())
		ERR(w, http.StatusInternalServerError, "internal server error")
	}
}

func stringNotEmpty(vs ...string) string {
	for i := 0; i < len(vs); i++ {
		if vs[i] != "" {
			return vs[i]
		}
	}
	return ""
}
