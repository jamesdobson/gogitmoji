package cmd

import (
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "🌍  Open gimoji information page in your browser",
	Long:  `Open gimoji information page in your browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		browser.OpenURL("https://gitmoji.carloscuesta.me/")
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
