package main

import (
	"strings"

	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/utils/fmtutil"
	"github.com/starudream/go-lib/selfupdate/v2"
)

var updateCmd = cobra.NewCommand(func(c *cobra.Command) {
	c.Use = "update"
	c.Short = "Update self"
	c.RunE = func(cmd *cobra.Command, args []string) error {
		source := &selfupdate.GoReleaser{
			Owner: "starudream",
			Repo:  "secret-tunnel",
			Name:  "stc",
		}
		confirm := func() bool {
			return strings.ToLower(strings.TrimSpace(fmtutil.Scan("update? (y/n):"))) == "y"
		}
		return selfupdate.Update(source, confirm)
	}
})

func init() {
	rootCmd.AddCommand(updateCmd)
}
