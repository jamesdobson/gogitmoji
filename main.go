package main

import (
	"fmt"

	"github.com/jamesdobson/gogitmoji/cmd"
	"github.com/spf13/cobra"
)

var (
	version = "local-dev"
	commit  = "no-commit-hash"
	date    = "unknown"
)

func main() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "ℹ️  Display the version of this program",
		Long:  `Display the version of this program.`,
		Run: func(*cobra.Command, []string) {
			fmt.Printf("This is gogitmoji %v (%v), build date: %v.\n\n", version, commit, date)
		},
	}
	cmd.AddCommand(versionCmd)

	// Transfer control to Cobra
	cmd.Execute()
}
