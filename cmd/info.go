package cmd

import (
	"log"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var infoURL = "https://gitmoji.carloscuesta.me/"

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "🌍  Open gimoji information page in your browser",
	Long:  `Open gimoji information page in your browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := browser.OpenURL(infoURL)

		if err != nil {
			log.Fatalf("Unable to launch page '%s': %v", infoURL, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
