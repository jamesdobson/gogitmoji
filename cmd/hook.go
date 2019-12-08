package cmd

import (
	"github.com/spf13/cobra"
)

// hookCmd represents the hook command
var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "🎣  Manage commit hooks",
	Long:  `Manage commit hooks.`,
}

func init() {
	rootCmd.AddCommand(hookCmd)
}
