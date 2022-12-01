package constant

import (
	"math"
	"sync"
	"time"

	"github.com/starudream/go-lib/config"
	consts "github.com/starudream/go-lib/constant"
)

var (
	VERSION = consts.VERSION
	BIDTIME = consts.BIDTIME
)

const (
	Name       = "SecretTunnel"
	DarwinName = "cn.starudream." + Name
	GitHub     = "https://github.com/starudream/secret-tunnel"

	MessageSize = math.MaxInt16

	ReadTimeout = 10 * time.Second
)

var (
	_init sync.Once

	_debug bool
)

func Debug() bool {
	_init.Do(func() { _debug = config.GetBool("debug") })
	return _debug
}
