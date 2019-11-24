package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available gitmoji",
	Long: `List all available gitmoji.

The gitmoji are printed on standard output, one gitmoji per line. Each line
has the emoji itself, the emoji's code, and a description of when to use
it.`,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list() {
	cache, err := NewGitmojiCache()

	if err != nil {
		log.Panic(err)
	}

	gitmojiList, err := cache.GetGitmoji()

	if err != nil {
		log.Panic("Unable to get list of gitmoji: ", err)
	}

	cyan := color.New(color.FgCyan)

	for i := 0; i < len(gitmojiList); i++ {
		gitmoji := gitmojiList[i]
		fmt.Printf("%s  - ", gitmoji.Emoji)
		cyan.Printf("%s", gitmoji.Code)
		fmt.Printf(" %s\n", gitmoji.Description)
	}

	fmt.Println("")
}
