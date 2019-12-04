package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/jamesdobson/gogitmoji/gitmoji"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "🔄  Update the list of gitmoji",
	Long: `Update the list of gitmoji.

Downloads a new list of gitmoji from https://gitmoji.carloscuesta.me/.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := gitmoji.UpdateCache()

		if err != nil {
			log.Fatalf("Unable to update: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
