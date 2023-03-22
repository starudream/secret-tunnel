package main

import (
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/flag"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/client"

	"github.com/starudream/secret-tunnel/internal/osx"
	"github.com/starudream/secret-tunnel/internal/service"
)

var (
	serviceCmd = &flag.Command{
		Use:   "service",
		Short: "Run as a service",
		Run: func(cmd *flag.Command, args []string) {
			svc := service.Get(client.Service)
			es := make(chan error, 100)
			go func() {
				for {
					e := <-es
					if e != nil {
						log.Warn().Msg(e.Error())
					}
				}
			}()
			_, _ = svc.Logger(es)
			osx.PA(svc.Run())
		},
	}

	serviceStatusCmd = &flag.Command{
		Use:   "status",
		Short: "Get the service status",
		Run: func(cmd *flag.Command, args []string) {
			st, se := service.Get(client.Service).Status()
			osx.PA(se, "the service status is "+service.StatusString(st))
		},
	}

	serviceStartCmd = &flag.Command{
		Use:   "start",
		Short: "Start the service",
		Run: func(cmd *flag.Command, args []string) {
			osx.PA(service.Get(client.Service).Start(), "the service is started")
		},
	}

	serviceStopCmd = &flag.Command{
		Use:   "stop",
		Short: "Stop the service",
		Run: func(cmd *flag.Command, args []string) {
			osx.PA(service.Get(client.Service).Stop(), "the service is stopped")
		},
	}

	serviceRestartCmd = &flag.Command{
		Use:   "restart",
		Short: "Restart the service",
		Run: func(cmd *flag.Command, args []string) {
			osx.PA(service.Get(client.Service).Restart(), "the service is started")
		},
	}

	serviceInstallCmd = &flag.Command{
		Use:   "install",
		Short: "Install the service",
		Run: func(cmd *flag.Command, args []string) {
			osx.PA(service.Get(client.Service).Install(), "the service is installed")
		},
	}

	serviceUninstallCmd = &flag.Command{
		Use:   "uninstall",
		Short: "Uninstall the service",
		Run: func(cmd *flag.Command, args []string) {
			osx.PA(service.Get(client.Service).Uninstall(), "the service is uninstalled")
		},
	}

	serviceReinstallCmd = &flag.Command{
		Use:   "reinstall",
		Short: "Reinstall the service",
		Run: func(cmd *flag.Command, args []string) {
			svc := service.Get(client.Service)
			osx.PE(svc.Uninstall())
			osx.PA(svc.Install(), "the service is installed")
		},
	}
)

func init() {
	serviceCmd.PersistentFlags().Bool("user", false, "run as current user, not root")
	osx.PE(config.BindPFlag("user", serviceCmd.PersistentFlags().Lookup("user")))

	serviceCmd.AddCommand(serviceStatusCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceRestartCmd)
	serviceCmd.AddCommand(serviceInstallCmd)
	serviceCmd.AddCommand(serviceUninstallCmd)
	serviceCmd.AddCommand(serviceReinstallCmd)
}
