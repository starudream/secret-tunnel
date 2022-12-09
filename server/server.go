package server

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/errx"
	"github.com/starudream/go-lib/log"
	"github.com/starudream/go-lib/seq"

	"github.com/starudream/secret-tunnel/constant"
	"github.com/starudream/secret-tunnel/internal/netx"
	"github.com/starudream/secret-tunnel/message"
	"github.com/starudream/secret-tunnel/model"
)

var COMM chan any

func Start(ctx context.Context) error {
	s, err := newServer(ctx)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	log.Info().Msgf("server start success on %s", s.address)

	s.ln = ln

	go s.comm()

	s.listener(s.ln)

	return nil
}

type Server struct {
	address string

	ctx    context.Context
	cancel context.CancelFunc

	ln net.Listener

	clients map[uint]string   // k: id, v: wid
	works   map[string]*iWork // k: wid
	cs      map[string]*iConn // key: sid
	workMu  sync.Mutex
}

type iWork struct {
	id string

	client *model.Client

	conn net.Conn
	c    *netx.Conn
}

type iConn struct {
	sid string
	tid string

	task *model.Task

	conn net.Conn
}

func newServer(ctx context.Context) (*Server, error) {
	ctx, cancel := context.WithCancel(ctx)
	s := &Server{
		address: config.GetString("addr"),
		ctx:     ctx,
		cancel:  cancel,
		clients: map[uint]string{},
		works:   map[string]*iWork{},
		cs:      map[string]*iConn{},
		workMu:  sync.Mutex{},
	}
	return s, nil
}

func (s *Server) comm() {
	for {
		v, ok := <-COMM
		if !ok {
			return
		}
		switch x := v.(type) {
		case *message.StopService:
			w := s.getWork(x.Cid)
			if w == nil {
				continue
			}
			message.WriteL(w.conn, &message.StopServiceReq{})
		case *message.UninstallService:
			w := s.getWork(x.Cid)
			if w == nil {
				continue
			}
			message.WriteL(w.conn, &message.UninstallServiceReq{})
		}
	}
}

func (s *Server) listener(ln net.Listener) {
	defer s.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Warn().Msgf("accept connection error: %v", err)
			netx.Close(conn)
			return
		}
		go s.into(conn)
	}
}

func (s *Server) into(conn net.Conn) {
	remote := netx.RemoteAddrString(conn)

	client, ok := s.login(conn)
	if !ok {
		log.Info().Str("remote", remote).Msgf("client login error")
		return
	}

	s.workMu.Lock()
	wid, exist := s.clients[client.Id]
	s.workMu.Unlock()

	if exist {
		log.Warn().Str("client", client.Name).Msgf("client has been registered, new replaces old")
		message.WriteL(s.works[wid].conn, &message.Close{Reason: "new client registered"})
		s.closeWork(wid)
	}

	log.Info().Str("client", client.Name).Msgf("client login success")

	c := netx.New(conn, true)

	for {
		sc, err := c.Session().Accept()
		if err != nil {
			if !errx.Is(err, io.EOF) && !strings.Contains(err.Error(), "use of closed network connection") {
				log.Warn().Str("client", client.Name).Msgf("accept connection error: %v", err)
			}
			log.Warn().Str("client", client.Name).Msgf("client disconnected")
			c.Close()
			go func() {
				ue := model.UpdateClientOffline(client.Id)
				if ue != nil {
					log.Warn().Str("client", client.Name).Msgf("update client offline error: %v", ue)
				}
			}()
			return
		}
		w := &iWork{
			id:     seq.NextId(),
			client: client,
			conn:   sc,
			c:      c,
		}
		go s.work(w)
	}
}

func (s *Server) login(conn net.Conn) (*model.Client, bool) {
	netx.SetReadTimeout(conn, constant.ReadTimeout)

	v, ok := message.ReadL(conn)
	if !ok {
		return nil, false
	}

	netx.SetReadTimeout(conn)

	loginReq, ok := v.(*message.LoginReq)
	if !ok {
		return nil, false
	}

	if loginReq.Key != "" {
		client, err := model.GetClientByKey(loginReq.Key)
		if err == nil {
			if !client.Active {
				log.Warn().Str("client", client.Name).Msgf("client not active")
				return nil, false
			}
			if message.WriteL(conn, &message.LoginResp{Ver: constant.VERSION}) {
				go func() {
					ue := model.UpdateClientOnline(&model.Client{
						Id:       client.Id,
						Addr:     netx.RemoteAddrString(conn),
						GO:       loginReq.GO,
						OS:       loginReq.OS,
						ARCH:     loginReq.ARCH,
						Hostname: loginReq.Hostname,
					})
					if ue != nil {
						log.Warn().Str("client", client.Name).Msgf("update client online error: %v", ue)
					}
				}()
				return client, true
			}
			return nil, false
		}
		log.Warn().Msgf("db get client error, %v", err)
	}

	message.WriteL(conn, &message.LoginResp{Ver: constant.VERSION, Error: message.NewError("key not match")})
	return nil, false
}

func (s *Server) work(w *iWork) (no bool) {
	defer func() {
		if !no {
			s.closeWork(w.id)
		}
	}()

	for {
		v, ok := message.ReadL(w.conn)
		if !ok {
			return
		}

		switch x := v.(type) {
		case *message.WorkReq:
			s.workMu.Lock()
			s.works[w.id] = w
			s.clients[w.client.Id] = w.id
			s.workMu.Unlock()
			if !message.WriteL(w.conn, &message.WorkResp{Wid: w.id}) {
				return
			}
			log.Debug().Str("client", w.client.Name).Msgf("work start")
		case *message.CreateTaskResp:
			if x.GetError() != "" {
				log.Warn().Str("client", w.client.Name).Msgf("create task error: %s", x.GetError())
				continue
			}
			s.workMu.Lock()
			c := s.cs[x.Sid]
			s.workMu.Unlock()
			if c.conn == nil {
				if c.tid != "" {
					message.WriteL(w.conn, &message.CloseTaskReq{Tid: c.tid})
				}
				continue
			}
			log.Debug().Str("client", w.client.Name).Str("task", c.task.Name).Msgf("task new connection")
			netx.Copy(w.conn, c.conn, "client", w.client.Name, "task", c.task.Name)
			s.workMu.Lock()
			delete(s.cs, x.Sid)
			s.workMu.Unlock()
			return
		case *message.ConnectTaskReq:
			t, err := model.GetTaskBySecret(0, x.Secret)
			if err != nil {
				log.Warn().Str("client", w.client.Name).Msgf("db get task error, %v", err)
				message.WriteL(w.conn, &message.ConnectTaskResp{Error: message.NewError("secret not match")})
				continue
			}
			if !t.Active {
				log.Warn().Str("client", w.client.Name).Str("task", t.Name).Msgf("task not active")
				message.WriteL(w.conn, &message.ConnectTaskResp{Error: message.NewError("task not active")})
				continue
			}
			sid, task := seq.NextId(), message.NewTask(t.Id, t.Name, t.Addr, t.Compress)
			if !s.createTask(t.ClientId, x.Tid, sid, task) {
				log.Warn().Str("client", w.client.Name).Str("task", t.Name).Msgf("create task error, target client is not online")
				message.WriteL(w.conn, &message.ConnectTaskResp{Error: message.NewError("target client is not online")})
				continue
			}
			if !message.WriteL(w.conn, &message.ConnectTaskResp{Sid: sid, Task: task}) {
				continue
			}
			s.workMu.Lock()
			s.cs[sid] = &iConn{sid: sid, tid: x.Tid, task: t, conn: w.conn}
			s.workMu.Unlock()
			return true
		}
	}
}

func (s *Server) closeWork(wid string) {
	s.workMu.Lock()
	if w, exist := s.works[wid]; exist {
		w.c.Close()
		delete(s.works, wid)
		delete(s.clients, w.client.Id)
	}
	s.workMu.Unlock()
}

func (s *Server) createTask(cid uint, tid, sid string, task message.Task) bool {
	w := s.getWork(cid)
	if w == nil {
		return false
	}
	return message.WriteL(w.conn, &message.CreateTaskReq{Tid: tid, Sid: sid, Task: task})
}

func (s *Server) getWork(cid uint) (w *iWork) {
	s.workMu.Lock()
	w = s.works[s.clients[cid]]
	s.workMu.Unlock()
	return
}

func (s *Server) Close() {
	_ = s.ln.Close()

	s.workMu.Lock()
	for _, w := range s.works {
		w.c.Close()
	}
	s.workMu.Unlock()

	s.cancel()
}
