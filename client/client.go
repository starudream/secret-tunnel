package client

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/starudream/go-lib/core/v2/config"
	"github.com/starudream/go-lib/core/v2/config/version"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/maputil"
	"github.com/starudream/go-lib/core/v2/utils/signalutil"
	"github.com/starudream/go-lib/core/v2/utils/timeutil"

	"github.com/starudream/secret-tunnel/message"
	"github.com/starudream/secret-tunnel/util"
)

type Client struct {
	addr string
	dns  string
	key  string

	tasks []*Task

	dialer net.Dialer
	conn   net.Conn
	mux    *message.MuxConn
	wid    atomic.Value
	lns    maputil.SyncMap[string, net.Listener] // local task id

	redo  chan struct{}
	notre atomic.Bool
}

type Task struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address"`
	Secret  string `json:"secret"`
}

func Run() error {
	c := &Client{
		addr: config.Get("addr").String(),
		dns:  config.Get("dns").String(),
		key:  config.Get("key").String(),
		redo: make(chan struct{}),
	}

	err := c.init()
	if err != nil {
		return err
	}

	err = c.connect()
	if err != nil {
		return err
	}

	go c.keepalive()

	<-signalutil.Defer(c.exit).Done()

	return nil
}

func (c *Client) init() error {
	var tasks []*Task

	err := config.Unmarshal("tasks", &tasks)
	if err != nil {
		return err
	}

	for k, v := range config.Raw() {
		if !strings.HasPrefix(k, "task") {
			continue
		}
		if i, _ := strconv.Atoi(k[4:]); i <= 0 {
			continue
		}
		if m, ok := v.(map[string]any); ok {
			task := &Task{}
			task.Name, _ = m["name"].(string)
			task.Address, _ = m["address"].(string)
			task.Secret, _ = m["secret"].(string)
			tasks = append(tasks, task)
		}
	}

	for i, task := range tasks {
		if task.Secret == "" {
			slog.Warn("task %d secret is empty, ignore", i)
			continue
		}
		h, p, e := net.SplitHostPort(task.Address)
		if e != nil {
			slog.Warn("task %d address %s invalid, skip", i, task.Address)
			continue
		}
		c.tasks = append(c.tasks, &Task{
			Id:      util.UUIDShort(),
			Name:    task.Name,
			Address: net.JoinHostPort(h, p),
			Secret:  task.Secret,
		})
	}

	c.dialer = net.Dialer{Timeout: 10 * time.Second}
	if c.dns != "" {
		h, p, e := net.SplitHostPort(c.dns)
		if e != nil {
			slog.Warn("dns %s invalid, ignore", c.dns)
		} else {
			if p == "" {
				p = "53"
			}
			c.dns = net.JoinHostPort(h, p)
			slog.Debug("use dns %s", c.dns)
			c.dialer.Resolver = &net.Resolver{PreferGo: true}
			c.dialer.Resolver.Dial = func(context.Context, string, string) (net.Conn, error) {
				return (&net.Dialer{Timeout: 3 * time.Second}).Dial("udp", c.dns)
			}
		}
	}

	return nil
}

func (c *Client) connect() (err error) {
	c.conn, err = c.dialer.Dial("tcp", c.addr)
	if err != nil {
		return
	}

	err = c.login(message.NewConn(c.conn))
	if err != nil {
		return
	}

	c.mux = message.NewMuxConn(c.conn, false)

	go c.work()

	return
}

func (c *Client) keepalive() {
	attempt := 0

	for {
		<-c.redo

		if c.notre.Load() {
			return
		}

		c.close()

		next := time.Hour

		if attempt == 0 {
			slog.Warn("lost connection, reconnecting...")
		}

		if attempt < 999 {
			next = timeutil.JitterDuration(time.Second, 10*time.Second, attempt)
		} else if attempt < 99999 {
			next = timeutil.JitterDuration(time.Second, 10*time.Minute, attempt)
		}

		next = next.Truncate(time.Millisecond)

		slog.Info("reconnect after %s", next, slog.Int("attempt", attempt))

		<-time.After(next)

		attempt++

		err := c.connect()
		if err != nil {
			slog.Warn("reconnect error: %v", err, slog.Int("attempt", attempt))
			continue
		}

		slog.Info("reconnect success", slog.Int("attempt", attempt))

		c.redo = make(chan struct{})

		attempt = 0
	}
}

func (c *Client) login(conn *message.Conn) error {
	loginReq := &message.LoginReq{
		Ver:      version.GetVersionInfo().GitVersion,
		Key:      c.key,
		GO:       runtime.Version(),
		OS:       runtime.GOOS,
		ARCH:     runtime.GOARCH,
		Hostname: func() string { s, _ := os.Hostname(); return s }(),
	}
	if !conn.WriteMessage(loginReq) {
		return fmt.Errorf("write login request error")
	}

	v, err := conn.ReadMessage(10 * time.Second)
	if err != nil {
		return err
	}

	loginResp, ok := v.(*message.LoginResp)
	if !ok {
		return fmt.Errorf("invalid login response")
	}

	if msg := loginResp.GetError(); msg != "" {
		return fmt.Errorf(msg)
	}

	return nil
}

func (c *Client) work() {
	mux, err := c.mux.Open()
	if err != nil {
		slog.Warn("open mux error: %v", err)
		return
	}

	defer close(c.redo)

	if !mux.WriteMessage(&message.WorkReq{}) {
		return
	}

	for {
		v, e := mux.ReadMessage()
		if e != nil {
			if message.ErrOther(err) {
				slog.Warn("read message error: %v", e)
			}
			return
		}

		switch x := v.(type) {
		case *message.Close:
			slog.Warn("closed by server: %s", x.Reason)
			signalutil.Cancel()
			return
		case *message.WorkResp:
			if msg := x.GetError(); msg != "" {
				slog.Error("register as work error: %s", msg)
				return
			}
			c.wid.Store(x.Wid)
			slog.Info("register as work success")
			for _, task := range c.tasks {
				go func(task *Task) { c.handleTask(task) }(task)
			}
		case *message.CreateTaskReq:
			go c.createTask(x)
		}
	}
}

func (c *Client) handleTask(task *Task) {
	c.closeTask(task.Id)

	ln, err := net.Listen("tcp", task.Address)
	if err != nil {
		slog.Error("listen to local %s error: %v", task.Address, err)
		return
	}

	c.lns.Store(task.Id, ln)
	defer c.closeTask(task.Id)

	str := task.Address
	if task.Name != "" {
		str = task.Name + " " + str
	}
	slog.Info("listen local task %s", str)

	for {
		local, e := message.AcceptConn(ln)
		if e != nil {
			if message.ErrOther(e) {
				slog.Warn("accept local %s error: %v", task.Address, e)
			}
			gh.Close(local)
			return
		}
		go c.connectTask(local, task)
	}
}

func (c *Client) connectTask(local *message.Conn, task *Task) {
	remote, err := c.mux.Open()
	if err != nil {
		slog.Warn("open mux error: %v", err)
		return
	}
	defer gh.Close(remote)

	connectReq := &message.ConnectTaskReq{
		Wid:  c.wid.Load().(string),
		Tid:  task.Id,
		Task: message.Task{Secret: task.Secret},
	}
	if !remote.WriteMessage(connectReq) {
		return
	}

	v, err := remote.ReadMessage(10 * time.Second)
	if err != nil {
		return
	}

	connectResp, ok := v.(*message.ConnectTaskResp)
	if !ok {
		return
	}

	if msg := connectResp.GetError(); msg != "" {
		slog.Error("task %s connect error: %s", task.Name, msg)
		return
	}

	slog.Debug("task %s new connection", connectResp.Task.Name, slog.String("sid", connectResp.Sid))

	if connectResp.Compress {
		message.Copy(message.WithSnappy(remote), local)
	} else {
		message.Copy(remote, local)
	}
}

func (c *Client) createTask(x *message.CreateTaskReq) {
	remote, err := c.mux.Open()
	if err != nil {
		slog.Warn("open mux error: %v", err)
		return
	}
	defer gh.Close(remote)

	local, err := net.DialTimeout("tcp", x.Task.Addr, 10*time.Second)
	if err != nil {
		slog.Error("dial to local %s error: %v", x.Task.Addr, err)
		remote.WriteMessage(&message.CreateTaskResp{Sid: x.Sid, Error: message.NewError("dial error")})
		return
	}
	defer gh.Close(local)

	if !remote.WriteMessage(&message.CreateTaskResp{Sid: x.Sid}) {
		return
	}

	slog.Debug("task %s new connection", x.Task.Name, slog.String("sid", x.Sid))

	if x.Task.Compress {
		message.Copy(local, message.WithSnappy(remote))
	} else {
		message.Copy(local, remote)
	}
}

func (c *Client) closeTask(id string) {
	ln, exists := c.lns.LoadAndDelete(id)
	if exists {
		gh.Close(ln)
	}
}

func (c *Client) exit() {
	c.notre.Store(true)
	c.close()
}

func (c *Client) close() {
	c.lns.Range(func(key string, ln net.Listener) bool {
		gh.Close(ln)
		c.lns.Delete(key)
		return true
	})

	gh.Close(c.mux)
	c.mux = nil
	gh.Close(c.conn)
	c.conn = nil

	c.wid.Store("")
}
