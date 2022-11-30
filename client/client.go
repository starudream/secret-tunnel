package client

import (
	"context"
	"fmt"
	"math"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/starudream/go-lib/codec/json"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/log"
	"github.com/starudream/go-lib/seq"

	"github.com/starudream/secret-tunnel/constant"
	"github.com/starudream/secret-tunnel/message"
	"github.com/starudream/secret-tunnel/netx"
)

func Start(ctx context.Context) error {
	c, err := newClient(ctx)
	if err != nil {
		return err
	}

	err = c.init()
	if err != nil {
		return err
	}

	go c.keep()

	<-c.ctx.Done()

	log.Info().Msgf("client exit")

	return nil
}

type Client struct {
	dns     string
	address string
	key     string
	tasks   []*iTask

	dialer net.Dialer

	exit   uint32
	ctx    context.Context
	cancel context.CancelFunc
	lostCh chan struct{}

	conn net.Conn
	c    *netx.Conn
	cMu  sync.Mutex

	wid    string
	lns    map[string]net.Listener // k: sid
	workMu sync.Mutex
}

type iTask struct {
	Id      string `json:"id,omitempty"`
	Address string `json:"address"`
	Secret  string `json:"secret"`
}

func newClient(ctx context.Context) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := &Client{
		dns:     config.GetString("dns"),
		address: config.GetString("addr"),
		key:     config.GetString("key"),
		ctx:     ctx,
		cancel:  cancel,
		lostCh:  make(chan struct{}),
		cMu:     sync.Mutex{},
		lns:     map[string]net.Listener{},
		workMu:  sync.Mutex{},
	}
	tasks, err := json.UnmarshalTo[[]*iTask]([]byte(config.GetString("tasks")))
	if err != nil {
		return nil, err
	}
	if len(tasks) > 0 {
		for i := 0; i < len(tasks); i++ {
			t := tasks[i]
			if t.Secret == "" {
				continue
			}
			h, p, e := net.SplitHostPort(t.Address)
			if e != nil {
				log.Warn().Msgf("invalid address: %s", t.Address)
				continue
			}
			c.tasks = append(c.tasks, &iTask{Id: seq.NextId(), Address: net.JoinHostPort(h, p), Secret: t.Secret})
		}
	}
	if c.dns != "" {
		if !strings.Contains(c.dns, ":") {
			c.dns += ":53"
		}
		c.dialer = net.Dialer{
			Timeout: 10 * time.Second,
			Resolver: &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					return net.Dial("udp", c.dns)
				},
			},
		}
	}
	return c, nil
}

func (c *Client) init() (err error) {
	err = c.dial()
	if err != nil {
		log.Warn().Msgf("dial to server error: %v", err)
		c.c.Close()
		return err
	}

	err = c.login()
	if err != nil {
		log.Warn().Msgf("login to server error: %v", err)
		c.c.Close()
		return err
	}

	c.cMu.Lock()
	c.c = netx.New(c.conn, false)
	c.cMu.Unlock()

	log.Info().Msgf("login to server success")

	go c.work()

	return nil
}

func (c *Client) dial() error {
	conn, err := c.dialer.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	c.cMu.Lock()
	c.conn = conn
	c.cMu.Unlock()

	return nil
}

func (c *Client) login() error {
	loginReq := &message.LoginReq{
		Ver:      constant.VERSION,
		Key:      c.key,
		GO:       runtime.Version(),
		OS:       runtime.GOOS,
		ARCH:     runtime.GOARCH,
		Hostname: func() string { name, _ := os.Hostname(); return name }(),
	}

	err := message.Write(c.conn, loginReq)
	if err != nil {
		return err
	}

	netx.SetReadTimeout(c.conn, constant.ReadTimeout)

	v, err := message.Read(c.conn)
	if err != nil {
		return err
	}

	netx.SetReadTimeout(c.conn)

	loginResp := v.(*message.LoginResp)

	if loginResp.GetError() != "" {
		return fmt.Errorf("%s", loginResp.GetError())
	}

	return nil
}

func (c *Client) keep() {
	count := 1

	for {
		<-c.lostCh

		if atomic.LoadUint32(&c.exit) != 0 {
			return
		}

		duration := func() time.Duration {
			if count <= 5 {
				return time.Duration(math.Pow(2, float64(count))) * time.Second
			} else {
				return time.Minute
			}
		}()

		if count == 1 {
			log.Warn().Msgf("lost connection with server")
		}

		time.Sleep(duration)
		count++

		log.Info().Msgf("try to reconnect to server")

		err := c.init()
		if err != nil {
			continue
		}

		c.lostCh = make(chan struct{})

		count = 1
	}
}

func (c *Client) work() {
	conn, err := c.c.Session().Open()
	if err != nil {
		log.Warn().Msgf("open stream error: %v", err)
		c.c.Close()
		return
	}

	defer close(c.lostCh)

	if !message.WriteL(conn, &message.WorkReq{}) {
		c.c.Close()
		return
	}

	for {
		v, ok := message.ReadL(conn)
		if !ok {
			c.c.Close()
			return
		}

		switch x := v.(type) {
		case *message.Close:
			log.Warn().Msgf("closed by server: %s", x.Reason)
			c.Close()
			return
		case *message.WorkResp:
			if x.GetError() != "" {
				log.Warn().Msgf("register as work error: %v", x.GetError())
				c.c.Close()
				return
			}
			c.wid = x.Wid
			for _, t := range c.tasks {
				t := t
				go c.connectTask(t)
			}
		case *message.CreateTaskReq:
			go c.createTask(x.Sid, x.Task)
		case *message.CloseTaskReq:
			go c.closeTask(x.Tid)
		}
	}
}

func (c *Client) connectTask(t *iTask) {
	ln, err := net.Listen("tcp", t.Address)
	if err != nil {
		log.Warn().Msgf("listen to local %s error: %v", t.Address, err)
		return
	}
	defer c.closeTask(t.Id)

	c.workMu.Lock()
	c.lns[t.Id] = ln
	c.workMu.Unlock()

	for {
		local, ae := ln.Accept()
		if ae != nil {
			return
		}
		go c.copyConn(local, t)
	}
}

func (c *Client) copyConn(local net.Conn, t *iTask) {
	remote, err := c.c.Session().Open()
	if err != nil {
		log.Warn().Msgf("open stream error: %v", err)
		return
	}
	defer netx.Close(remote)

	task := message.NewTaskSecret(t.Secret)

	if !message.WriteL(remote, &message.ConnectTaskReq{Wid: c.wid, Tid: t.Id, Task: task}) {
		return
	}

	netx.SetReadTimeout(remote, constant.ReadTimeout)

	v, ok := message.ReadL(remote)
	if !ok {
		return
	}

	netx.SetReadTimeout(remote)

	connectResp := v.(*message.ConnectTaskResp)

	log.Debug().Str("task", connectResp.Task.Name).Str("sid", connectResp.Sid).Str("tid", t.Id).Msgf("task new connection")

	netx.Copy(remote, local)
}

func (c *Client) createTask(sid string, t message.Task) {
	remote, err := c.c.Session().Open()
	if err != nil {
		log.Warn().Msgf("open stream error: %v", err)
		c.c.Close()
		return
	}
	defer netx.Close(remote)

	local, err := net.DialTimeout("tcp", t.Addr, 10*time.Second)
	if err != nil {
		log.Warn().Msgf("dial to local %s error: %v", t.Addr, err)
		message.WriteL(remote, &message.CreateTaskResp{Sid: sid, Error: message.NewError("dial to local %s error: %v", t.Addr, err)})
		return
	}
	defer netx.Close(local)

	if !message.WriteL(remote, &message.CreateTaskResp{Sid: sid}) {
		return
	}

	log.Debug().Str("sid", sid).Str("addr", t.Addr).Msgf("task new connection")

	netx.Copy(local, remote)
}

func (c *Client) closeTask(tid string) {
	c.workMu.Lock()
	if ln, exist := c.lns[tid]; exist {
		_ = ln.Close()
	}
	delete(c.lns, tid)
	c.workMu.Unlock()
}

func (c *Client) Close() {
	atomic.StoreUint32(&c.exit, 1)

	c.c.Close()

	c.cancel()
}
