package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/core/v2/config/version"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/maputil"
	"github.com/starudream/go-lib/core/v2/utils/signalutil"
	"github.com/starudream/go-lib/ntfy/v2"

	"github.com/starudream/secret-tunnel/message"
	"github.com/starudream/secret-tunnel/model"
	"github.com/starudream/secret-tunnel/util"
)

type Server struct {
	addr string
	ln   net.Listener

	clients  maputil.SyncMap[uint, *Work]
	works    maputil.SyncMap[string, *Work]
	sessions maputil.SyncMap[string, *Session]
}

type Work struct {
	id     string
	client *model.Client
	conn   *message.Conn
}

type Session struct {
	id   string
	task *model.Task
	conn *message.Conn
	wait chan struct{}
}

func Run() (err error) {
	s := &Server{
		addr: config.Get("addr").String(),
	}

	s.ln, err = net.Listen("tcp", s.addr)
	if err != nil {
		return
	}

	slog.Info("server start success on %s", s.ln.Addr().String())

	go s.accept()

	<-signalutil.Defer(s.close).Done()

	return
}

func (s *Server) accept() {
	defer gh.Close(s.ln)

	for {
		conn, err := message.AcceptConn(s.ln)
		if err != nil {
			if message.ErrOther(err) {
				slog.Warn("accept error: %v", err)
			}
			gh.Close(conn)
			return
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *message.Conn) {
	defer gh.Close(conn)

	client, err := s.clientLogin(conn)
	if err != nil {
		slog.Warn("login error: %v", err, slog.String("remote", conn.RemoteAddrString()))
		return
	}

	work, exists := s.clients.LoadAndDelete(client.Id)
	if exists {
		slog.Warn("client %s already exists, kick out the old one", client.Name,
			slog.String("client", client.Name),
			slog.String("old", work.conn.RemoteAddrString()), slog.String("new", conn.RemoteAddrString()),
		)
		work.conn.WriteMessage(&message.Close{Reason: "kick out"})
		s.closeWork(work.id)
	}

	s.acceptMux(client, message.NewMuxConn(conn.Conn, true))
}

func (s *Server) clientLogin(conn *message.Conn) (_ *model.Client, err error) {
	v, err := conn.ReadMessage(10 * time.Second)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			conn.WriteMessage(&message.LoginResp{Ver: version.GetVersionInfo().GitVersion, Error: message.NewError(err.Error())})
		}
	}()

	loginReq, ok := v.(*message.LoginReq)
	if !ok {
		return nil, fmt.Errorf("invalid login request")
	}

	if loginReq.Key == "" {
		return nil, fmt.Errorf("invalid login key")
	}

	client, err := model.GetClientByKey(loginReq.Key)
	if err != nil {
		return nil, err
	}

	if !client.Active {
		return nil, fmt.Errorf("client not active")
	}

	if !conn.WriteMessage(&message.LoginResp{Ver: version.GetVersionInfo().GitVersion}) {
		return nil, fmt.Errorf("write login response error")
	}

	err = model.UpdateClientOnline(&model.Client{
		Id:       client.Id,
		Ver:      loginReq.Ver,
		Addr:     conn.RemoteAddrString(),
		GO:       loginReq.GO,
		OS:       loginReq.OS,
		ARCH:     loginReq.ARCH,
		Hostname: loginReq.Hostname,
	})
	if err != nil {
		slog.Error("update client online error: %v", err)
	}

	return client, nil
}

func (s *Server) acceptMux(client *model.Client, mux *message.MuxConn) {
	defer func() {
		gh.Close(mux)
		slog.Info("client %s disconnected", client.Name)
		err := model.UpdateClientOffline(client.Id)
		if err != nil {
			slog.Error("update client offline error: %v", err)
		}
		err = ntfy.Notify(context.Background(), fmt.Sprintf("%s offline", client.Name))
		if err != nil && !errors.Is(err, ntfy.ErrNoConfig) {
			slog.Warn("notify offline error: %v", err)
		}
	}()

	go func() {
		slog.Info("client %s connected", client.Name)
		err := ntfy.Notify(context.Background(), fmt.Sprintf("%s online", client.Name))
		if err != nil && !errors.Is(err, ntfy.ErrNoConfig) {
			slog.Warn("notify online error: %v", err)
		}
	}()

	for {
		conn, err := message.AcceptConn(mux.Session())
		if err != nil {
			if message.ErrOther(err) {
				slog.Warn("accept mux error: %v", err)
			}
			return
		}
		go s.handleWork(&Work{
			id:     util.UUIDShort(),
			client: client,
			conn:   conn,
		})
	}
}

func (s *Server) handleWork(w *Work) (skipClose bool) {
	defer func() {
		if !skipClose {
			s.closeWork(w.id)
		}
	}()

	for {
		v, err := w.conn.ReadMessage()
		if err != nil {
			if message.ErrOther(err) {
				slog.Warn("read message error: %v", err)
			}
			return
		}

		switch x := v.(type) {
		case *message.WorkReq:
			s.works.Store(w.id, w)
			s.clients.Store(w.client.Id, w)
			if !w.conn.WriteMessage(&message.WorkResp{Wid: w.id}) {
				return
			}
			slog.Debug("client %s registered as work", w.client.Name)
		case *message.ConnectTaskReq:
			err = s.connectTask(w.conn, x)
			if err != nil {
				slog.Warn("connect task error: %v", err)
			} else {
				return true
			}
		case *message.CreateTaskResp:
			err = s.createTask(w.conn, x)
			if err != nil {
				slog.Warn("create task error: %v", err)
			} else {
				return
			}
		}
	}
}

func (s *Server) connectTask(conn *message.Conn, x *message.ConnectTaskReq) (err error) {
	defer func() {
		if err != nil {
			conn.WriteMessage(&message.ConnectTaskResp{Error: message.NewError(err.Error())})
		}
	}()

	task, err := model.GetTaskBySecret(0, x.Secret)
	if err != nil {
		return err
	}

	if !task.Active {
		return fmt.Errorf("task not active")
	}

	work, exists := s.clients.Load(task.ClientId)
	if !exists {
		return fmt.Errorf("target client is offline")
	}

	session := &Session{id: util.UUIDShort(), task: task, conn: conn, wait: make(chan struct{})}

	s.sessions.Store(session.id, session)

	taskMsg := message.Task{Id: task.Id, Name: task.Name, Addr: task.Addr, Compress: task.Compress}

	if !work.conn.WriteMessage(&message.CreateTaskReq{Tid: x.Tid, Sid: session.id, Task: taskMsg}) {
		return fmt.Errorf("write create task request error")
	}

	select {
	case <-time.After(3 * time.Second):
		return fmt.Errorf("create task timeout")
	case <-session.wait:
		// success
	}

	if !conn.WriteMessage(&message.ConnectTaskResp{Sid: session.id, Task: taskMsg}) {
		return fmt.Errorf("write connect task response error")
	}

	return nil
}

func (s *Server) createTask(conn *message.Conn, x *message.CreateTaskResp) (err error) {
	if msg := x.GetError(); msg != "" {
		return fmt.Errorf(msg)
	}

	session, exists := s.sessions.Load(x.Sid)
	if !exists {
		return fmt.Errorf("session not exists")
	}

	close(session.wait)

	in, out := message.Copy(conn, session.conn)

	slog.Debug("task %s traffic: %s in, %s out", session.task.Name, model.Size(in), model.Size(out), slog.String("sid", x.Sid))

	err = model.UpdateTaskTraffic(session.task.Id, uint(in), uint(out))
	if err != nil {
		slog.Error("update task traffic error: %v", err)
	}

	s.sessions.Delete(x.Sid)

	return nil
}

func (s *Server) closeWork(id string) {
	work, exists := s.works.LoadAndDelete(id)
	if exists {
		gh.Close(work.conn)
		s.clients.Delete(work.client.Id)
	}
}

func (s *Server) close() {
	s.works.Range(func(key string, work *Work) bool {
		gh.Close(work.conn)
		s.works.Delete(key)
		s.clients.Delete(work.client.Id)
		return true
	})

	s.sessions.Range(func(key string, session *Session) bool {
		gh.Close(session.conn)
		s.sessions.Delete(key)
		return true
	})

	gh.Close(s.ln)
	s.ln = nil
}
