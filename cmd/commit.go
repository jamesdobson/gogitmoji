package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jamesdobson/gogitmoji/gitmoji"
	"github.com/jamesdobson/gogitmoji/tmpl"
)

const (
	formatSetting = "format"
	formatAsEmoji = "emoji"

	scopeSetting = "scope"

	typeSetting = "type"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "⚡️  Compose a commit message and execute git commit",
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
	// TODO: Change to "template"
	commitCmd.Flags().StringP("type", "t", tmpl.DefaultTemplateName, `Commit template name.`)

	err = viper.BindPFlag(formatSetting, commitCmd.Flags().Lookup("format"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag(scopeSetting, commitCmd.Flags().Lookup("scope"))
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag(typeSetting, commitCmd.Flags().Lookup("type"))
	if err != nil {
		panic(err)
	}
}

func commit() {
	templates := viper.GetStringMap("templates")
	t := viper.GetString(typeSetting)

	// TODO: Update documentation with an example of setting a template
	tmpl.LoadTemplates(templates)

	if tpl, ok := tmpl.TemplateLookup[t]; ok {
		commitWithTemplate(tpl)
		fmt.Printf("\ngogitmoji done.\n")
	} else {
		log.Fatalf("Unknown commit type: \"%s\"\n", t)
	}
}

func commitWithTemplate(tpl tmpl.CommandTemplate) {
	var answers = map[string]interface{}{}

	for q := 0; q < len(tpl.Prompts); q++ {
		var question = tpl.Prompts[q]

		if question.EnableSetting != "" && !viper.GetBool(question.EnableSetting) {
			continue
		}

		switch question.PromptType {
		case "text":
			answer := promptOrCancel(question.Prompt, question.Mandatory)
			answers[question.ValueCode] = answer

		case "choice":
			answer := promptChoice(question)
			answers[question.ValueCode] = answer

		case "gitmoji":
			gitmoji, err := promptGitmoji()

			if err != nil {
				if err == promptui.ErrInterrupt {
					fmt.Println("Canceled.")
					os.Exit(1)
				}

				log.Panic("Couldn't pick a gitmoji: ", err)
			}

			answers[question.ValueCode] = gitmoji

		default:
			log.Fatalf("Unknown prompt type '%s'...\n", question.PromptType)
		}
	}

	args := generateArgs(&tpl, answers)

	displayCommand := getPrintableCommand(tpl.Command, args)
	fmt.Printf("Going to execute: %s\n\n", displayCommand)

	if !confirm("Execute") {
		fmt.Printf("Canceled.\n")
		return
	}

	run(tpl.Command, args)
}

func generateArgs(tpl *tmpl.CommandTemplate, answers map[string]interface{}) []string {
	var args = make([]string, 0, len(tpl.CommandArgs))
	var sb strings.Builder
	var functions = map[string]interface{}{
		"getString": viper.GetString,
	}

	for n := 0; n < len(tpl.CommandArgs); n++ {
		sb.Reset()
		t, err := template.New("arg").
			Funcs(functions).
			Parse(tpl.CommandArgs[n])

		if err != nil {
			panic(err)
		}

		err = t.Execute(&sb, answers)

		if err != nil {
			panic(err)
		}

		arg := sb.String()

		if arg != "" {
			args = append(args, arg)
		}
	}

	return args
}

func getPrintableCommand(name string, args []string) string {
	var sb *strings.Builder = &strings.Builder{}

	sb.WriteString(name)

	isPlain := regexp.MustCompile(`^[-.A-Za-z0-9]+$`).MatchString

	for n := 0; n < len(args); n++ {
		if isPlain(args[n]) {
			fmt.Fprintf(sb, " %s", args[n])
		} else {
			fmt.Fprintf(sb, " \"%s\"", strings.ReplaceAll(args[n], `"`, `\"`))
		}
	}

	return sb.String()
}

func run(name string, args []string) {
	fmt.Printf("Executing...\n")

	var cmd = exec.Command(name, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		code := cmd.ProcessState.ExitCode()
		fmt.Printf("\n'%s' exited with code: %d\n", name, code)
		os.Exit(code)
	}
}

func confirm(question string) bool {
	prompt := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
	}

	result, err := prompt.Run()

	return err == nil && strings.ToLower(result) == "y"
}

func promptOrCancel(question string, mandatory bool) string {
	s, err := prompt(question, mandatory)

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Canceled.")
			os.Exit(1)
		}

		log.Panic(err)
	}

	return s
}

func prompt(question string, mandatory bool) (string, error) {
	templates := &promptui.PromptTemplates{
		Success: `{{ "✔" | faint }} {{ . | faint }}{{ ":" | faint }} `,
	}

	validator := func(input string) error {
		if !mandatory || len(input) >= 1 {
			return nil
		}

		return fmt.Errorf("this is required")
	}

	prompt := promptui.Prompt{
		Label:     question,
		Validate:  validator,
		Templates: templates,

		// Disable the pointer
		Pointer: func(x []rune) []rune { return x },
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func promptGitmoji() (gitmoji.Gitmoji, error) {
	cache, err := gitmoji.NewCache()

	if err != nil {
		log.Panic(err)
	}

	glist, err := cache.GetGitmoji()

	if err != nil {
		log.Panic("Unable to get list of gitmoji: ", err)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ \"?\" | yellow }} {{ . }}",
		Active:   "‣ {{ .Emoji }}  - {{ .Code | cyan }} - {{ .Description }}",
		Inactive: "  {{ .Emoji }}  - {{ .Code | cyan }} - {{ .Description }}",
		Selected: `{{ "? Choose a gitmoji" | faint }} {{ .Emoji }}  - {{ .Description }}`,
		Details: `
--------- Gitmoji ----------
{{ "Name:" | faint }}	{{ .Emoji }} {{ .Name }}
{{ "Entity:" | faint }}	{{ .Entity }}
{{ "Code:" | faint }}	{{ .Code }}
{{ "Description:" | faint }}	{{ .Description }}`,
	}

	searcher := func(input string, index int) bool {
		gitmoji := glist[index]
		tosearch := gitmoji.Name + gitmoji.Code + gitmoji.Description

		// Normalize
		tosearch = strings.Replace(strings.ToLower(tosearch), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(tosearch, input)
	}

	prompt := promptui.Select{
		Label:     "Choose a gitmoji",
		Items:     glist,
		Templates: templates,
		Size:      12,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return gitmoji.Gitmoji{}, err
	}

	return glist[i], nil
}

func promptChoice(question tmpl.Prompt) string {
	templates := &promptui.SelectTemplates{
		Label: "{{ \"?\" | yellow }} {{ . }}",
		Active: "‣ {{ .Value }} 	{{ .Description }}",
		Inactive: "  {{ .Value }} 	{{ .Description }}",
		Selected: `{{ "? ` + question.Prompt + `" | faint }} {{ .Value }}  - {{ .Description }}`,
	}

	searcher := func(input string, index int) bool {
		t := question.Choices[index]
		tosearch := t.Value + t.Description

		// Normalize
		tosearch = strings.Replace(strings.ToLower(tosearch), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(tosearch, input)
	}

	prompt := promptui.Select{
		Label:     "Choose the type of commit:",
		Items:     question.Choices,
		Templates: templates,
		Size:      12,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Canceled.")
			os.Exit(1)
		}

		log.Panic(err)
	}

	return question.Choices[i].Value
}
