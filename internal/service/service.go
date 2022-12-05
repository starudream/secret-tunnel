package service

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/kardianos/service"

	"github.com/starudream/go-lib/config"

	"github.com/starudream/secret-tunnel/constant"
)

func Get(i service.Interface) service.Service {
	serviceKV := service.KeyValue{}
	serviceKV["UserService"] = config.GetBool("user")

	envVars := map[string]string{}
	envVars[constant.PREFIX+"_USER"] = config.GetString("user")

	if executable, err := os.Executable(); err == nil {
		envVars[constant.PREFIX+"_LOG_FILE_PATH"] = executable + ".log"
	} else {
		envVars[constant.PREFIX+"_LOG_FILE_PATH"] = filepath.Join(os.TempDir(), "stc.log")
	}

	serviceCfg := &service.Config{
		Name:        constant.Name,
		DisplayName: constant.Name + "Client",
		Description: constant.GitHub,
		Arguments:   []string{"service", "--addr", config.GetString("addr"), "--key", config.GetString("key"), "--dns", config.GetString("dns")},
		Option:      serviceKV,
		EnvVars:     envVars,
	}

	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "darwin" {
		serviceCfg.Name = constant.DarwinName
	}

	svc, _ := service.New(i, serviceCfg)

	return svc
}

func StatusString(st service.Status) string {
	switch st {
	case service.StatusRunning:
		return "running"
	case service.StatusStopped:
		return "stopped"
	default:
		return "unknown"
	}
}
