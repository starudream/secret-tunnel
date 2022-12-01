package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/config"

	"github.com/starudream/secret-tunnel/client"
	"github.com/starudream/secret-tunnel/constant"
)

var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Run as a service",
		Args:  cobra.MinimumNArgs(1),
	}

	serviceStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Get the service status",
		Run: func(cmd *cobra.Command, args []string) {
			st, se := getService().Status()
			if se != nil {
				p(se)
			}
			ss(st)
		},
	}

	serviceStartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the service",
		Run: func(cmd *cobra.Command, args []string) {
			p(getService().Start(), "the service is started")
		},
	}

	serviceStopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop the service",
		Run: func(cmd *cobra.Command, args []string) {
			p(getService().Stop(), "the service is stopped")
		},
	}

	serviceRestartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart the service",
		Run: func(cmd *cobra.Command, args []string) {
			p(getService().Restart(), "the service is started")
		},
	}

	serviceInstallCmd = &cobra.Command{
		Use:   "install",
		Short: "Install the service",
		Run: func(cmd *cobra.Command, args []string) {
			p(getService().Install(), "the service is installed")
		},
	}

	serviceUninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall the service",
		Run: func(cmd *cobra.Command, args []string) {
			p(getService().Uninstall(), "the service is uninstalled")
		},
	}

	serviceReinstallCmd = &cobra.Command{
		Use:   "reinstall",
		Short: "Reinstall the service",
		Run: func(cmd *cobra.Command, args []string) {
			svc := getService()
			_ = svc.Uninstall()
			p(svc.Install(), "the service is installed")
		},
	}
)

func init() {
	serviceCmd.AddCommand(serviceStatusCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceRestartCmd)
	serviceCmd.AddCommand(serviceInstallCmd)
	serviceCmd.AddCommand(serviceUninstallCmd)
	serviceCmd.AddCommand(serviceReinstallCmd)

	serviceCmd.PersistentFlags().Bool("user", false, "run as current user, not root")
	_ = config.BindPFlag("user", serviceCmd.PersistentFlags().Lookup("user"))
}

func getService() service.Service {
	serviceKV := service.KeyValue{}
	serviceKV["UserService"] = config.GetBool("user")

	serviceCfg := &service.Config{
		Name:        constant.Name,
		DisplayName: constant.Name + "Client",
		Description: constant.GitHub,
		Arguments:   os.Args[3:],
		Option:      serviceKV,
	}

	if isDarwin() {
		serviceCfg.Name = constant.DarwinName
	}

	svc, err := service.New(&iService{}, serviceCfg)
	if err != nil {
		p(err)
	}

	return svc
}

type iService struct {
}

var _ service.Interface = (*iService)(nil)

func (p *iService) Start(_ service.Service) error {
	return client.Start(context.Background())
}

func (p *iService) Stop(_ service.Service) error {
	app.Stop()
	return nil
}

//goland:noinspection ALL
func isDarwin() bool {
	return runtime.GOOS == "darwin"
}

func ss(st service.Status) {
	s := "the service status is "
	switch st {
	case service.StatusRunning:
		s += "running"
	case service.StatusStopped:
		s += "stopped"
	default:
		s += "unknown"
	}
	fmt.Println(s)
}

func p(v ...any) {
	if len(v) == 0 {
		os.Exit(0)
	}
	c, s, w := 0, "", os.Stdout
	nh := func() {
		if len(v) >= 2 {
			switch y := v[1].(type) {
			case string:
				v = v[1:]
				s = y
			}
		}
	}
	switch x := v[0].(type) {
	case string:
		s = x
	case error:
		if x != nil {
			c, s, w = 1, x.Error(), os.Stderr
			v = v[:1]
		} else {
			nh()
		}
	case nil:
		nh()
	default:
		c, s = 1, fmt.Sprint(x)
	}
	if len(v) >= 2 {
		s = fmt.Sprintf(s, v[1:]...)
	}
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	if s != "" {
		_, _ = fmt.Fprintf(w, s)
	}
	os.Exit(c)
}
