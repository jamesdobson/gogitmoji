package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jamesdobson/gogitmoji/tmpl"
)

const (
	formatSetting = "format"
	formatAsEmoji = "emoji"

	scopeSetting = "scope"

	templateSetting = "template"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "⚡️  Compose a commit message and execute git commit (default command)",
	Long: `Compose a commit message and execute git commit.

Prompts for the gitmoji to use, as well as the commit message itself. Once
all prompts are filled out, executes git commit.

This is the default command when no other command is specified to gogitmoji.`,
	Run: func(cmd *cobra.Command, args []string) {
		commit()
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	var err error

	commitCmd.Flags().StringP("format", "f", formatAsEmoji, `Emoji format; either "emoji" or "code".`)
	commitCmd.Flags().BoolP("scope", "p", false, "Enable scope prompt")
	commitCmd.Flags().StringP("template", "t", tmpl.DefaultTemplateName, `Commit template name.`)

	err = viper.BindPFlag(formatSetting, commitCmd.Flags().Lookup("format"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag(scopeSetting, commitCmd.Flags().Lookup("scope"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag(templateSetting, commitCmd.Flags().Lookup("template"))
	if err != nil {
		panic(err)
	}
}

func commit() {
	templates := viper.GetStringMap("templates")
	t := viper.GetString(templateSetting)

	tmpl.LoadTemplates(templates)

	if tpl, ok := tmpl.TemplateLookup[t]; ok {
		tmpl.RunTemplateCommand(tpl)
		fmt.Printf("\ngogitmoji done.\n")
	} else {
		log.Fatalf("Unknown commit template: \"%s\"\n", t)
	}
}
