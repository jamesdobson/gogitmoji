package cmd

import (
	"fmt"
	"log"
	"os"

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
	Use:     "commit",
	Aliases: []string{"hook"},
	Short:   "⚡️  Compose a commit message and execute git commit (default command)",
	Long: `Compose a commit message and execute git commit.

Prompts for the gitmoji to use, as well as the commit message itself. Once
all prompts are filled out, executes git commit.

This is the default command when no other command is specified to gogitmoji.

The hook do command has the same format, except that it takes one argument: the
path to a file containing the commit message. The hook do command is only
intended to be called by a commit hook; not directly from the CLI. There must
be a space between "hook" and "do."`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.CalledAs() == "hook" {
			if len(args) == 2 && args[0] == "do" {
				do(args[1:])
			} else {
				log.Fatalf("Argument to hook must be 'do', followed by path to commit message file.")
			}
		} else {
			commit()
		}
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

func do(args []string) {
	templates := viper.GetStringMap("templates")
	t := viper.GetString(templateSetting)

	tmpl.LoadTemplates(templates)

	if tpl, ok := tmpl.TemplateLookup[t]; ok {
		msg := tmpl.GetTemplateMessage(tpl)
		f, err := os.OpenFile(args[0], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

		if err != nil {
			log.Fatalf("Error opening commit message file: %v\n", err)
		}

		_, err = f.WriteString(msg)

		if err != nil {
			log.Fatalf("Error writing commit message file: %v\n", err)
		}
	} else {
		log.Fatalf("Unknown commit template: \"%s\"\n", t)
	}
}
