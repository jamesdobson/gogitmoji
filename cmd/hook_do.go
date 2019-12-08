package cmd

import (
	"log"
	"os"

	"github.com/jamesdobson/gogitmoji/tmpl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "❗️  Prompt the user and write message to provided file",
	Long: `Prompt the user and write message to provided file.

This command is only intended to be called by a commit hook.`,
	Run: func(cmd *cobra.Command, args []string) {
		do(args)
	},
}

func init() {
	hookCmd.AddCommand(doCmd)

	var err error

	doCmd.Flags().StringP("format", "f", formatAsEmoji, `Emoji format; either "emoji" or "code".`)
	doCmd.Flags().BoolP("scope", "p", false, "Enable scope prompt")
	doCmd.Flags().StringP("template", "t", tmpl.DefaultTemplateName, `Commit template name.`)

	err = viper.BindPFlag(formatSetting, doCmd.Flags().Lookup("format"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag(scopeSetting, doCmd.Flags().Lookup("scope"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag(templateSetting, doCmd.Flags().Lookup("template"))
	if err != nil {
		panic(err)
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
