package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "🔄  Update the list of gitmoji",
	Long: `Update the list of gitmoji.

Downloads a new list of gitmoji from https://gitmoji.carloscuesta.me/.`,
	Run: func(cmd *cobra.Command, args []string) {
		UpdateGitmojiCache()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
