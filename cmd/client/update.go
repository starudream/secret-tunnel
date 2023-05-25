package main

import (
	"github.com/starudream/go-lib/flag"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/secret-tunnel/constant"

	"github.com/starudream/secret-tunnel/internal/osx"
	"github.com/starudream/secret-tunnel/internal/update"
)

var (
	updateCmd = &flag.Command{
		Use:   "update",
		Short: "Update self",
		Run: func(cmd *flag.Command, args []string) {
			if updateTarget == "" {
				name, err := update.GetLatestReleaseName()
				osx.PE(err)

				updateTarget = name

				if updateTarget == constant.VERSION {
					osx.PA("already up-to-date")
				}

				log.Info().Msgf("latest version: %s, downloading...", updateTarget)

				file, err := update.DownloadReleaseFile(updateTarget, "", "", true)
				osx.PE(err)

				osx.PE(update.Apply(file))

				log.Info().Msgf("update to %s successfully", updateTarget)
			}
		},
	}

	updateTarget string
)

func init() {
	updateCmd.PersistentFlags().StringVar(&updateTarget, "target", "", "target version")
}
