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
)

const (
	formatSetting = "format"
	formatAsEmoji = "emoji"
	formatAsCode  = "code"

	scopeSetting = "scope"

	typeSetting      = "type"
	typeGitmoji      = "gitmoji"
	typeConventional = "conventional"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "⚡️  Compose a commit message and execute git commit",
	Long: `Compose a commit message and execute git commit.

Prompts for the gitmoji to use, as well as the commit message itself. Once
all prompts are filled out, executes git commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		commit()
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("format", "f", formatAsEmoji, `Emoji format; either "emoji" or "code".`)
	commitCmd.Flags().BoolP("scope", "p", false, "Enable scope prompt")
	commitCmd.Flags().StringP("type", "t", typeGitmoji, `Commit format; either "gitmoji" or "conventional".`)
	viper.BindPFlag(formatSetting, commitCmd.Flags().Lookup("format"))
	viper.BindPFlag(scopeSetting, commitCmd.Flags().Lookup("scope"))
	viper.BindPFlag(typeSetting, commitCmd.Flags().Lookup("type"))
}

func commit() {
	t := viper.GetString(typeSetting)
	var tpl CommitTemplate

	switch t {
	case typeGitmoji:
		tpl = GitmojiCommit
	case typeConventional:
		tpl = ConventionalCommit
	default:
		log.Fatalf("Unknown commit type: \"%s\"\n", t)
	}

	commitWithTemplate(tpl)

	fmt.Printf("\ngogitmoji done.\n")
}

func commitWithTemplate(tpl CommitTemplate) {
	var answers = map[string]interface{}{}

	for q := 0; q < len(tpl.Questions); q++ {
		var question = tpl.Questions[q]

		if question.EnableSetting != "" && !viper.GetBool(question.EnableSetting) {
			continue
		}

		switch question.PromptType {
		case "text":
			answer := promptOrCancel(question.Prompt, question.Mandatory)
			answers[question.ValueCode] = answer

		case "conventional":
			t, err := promptConventionalCommitType()

			if err != nil {
				if err == promptui.ErrInterrupt {
					fmt.Println("Canceled.")
					os.Exit(1)
				}

				log.Panic("Couldn't pick a conventional commit type: ", err)
			}

			answers[question.ValueCode] = t.Name

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

	var args = make([]string, 0, len(tpl.CommandArgs))
	var sb strings.Builder

	for n := 0; n < len(tpl.CommandArgs); n++ {
		sb.Reset()
		t, err := template.New("arg").
			Funcs(
				map[string]interface{}{
					"getString": func(name string) string {
						return viper.GetString(name)
					},
				}).
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

	fmt.Printf("Going to execute: %s", tpl.Command)

	isPlain := regexp.MustCompile(`^[-A-Za-z0-9]+$`).MatchString

	for n := 0; n < len(args); n++ {
		if isPlain(args[n]) {
			fmt.Printf(" %s", args[n])
		} else {
			fmt.Printf(" \"%s\"", args[n])
		}
	}

	fmt.Printf("\n\n")

	if !confirm("Execute") {
		fmt.Printf("Canceled.\n")
		return
	}

	run(tpl.Command, args)
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

type CommitQuestion struct {
	PromptType    string
	Mandatory     bool
	Prompt        string
	ValueCode     string
	EnableSetting string
}

type CommitTemplate struct {
	Questions   []CommitQuestion
	Command     string
	CommandArgs []string
}

var GitmojiCommit = CommitTemplate{
	Questions: []CommitQuestion{
		CommitQuestion{
			PromptType: "gitmoji",
			Mandatory:  true,
			ValueCode:  "gitmoji",
		},
		CommitQuestion{
			PromptType:    "text",
			Mandatory:     false,
			Prompt:        "Enter the scope of current changes",
			ValueCode:     "Scope",
			EnableSetting: scopeSetting,
		},
		CommitQuestion{
			PromptType: "text",
			Mandatory:  true,
			Prompt:     "Enter the commit title",
			ValueCode:  "title",
		},
		CommitQuestion{
			PromptType: "text",
			Mandatory:  false,
			Prompt:     "Enter the (optional) commit message",
			ValueCode:  "message",
		},
	},
	Command: "git",
	CommandArgs: []string{
		"commit",
		"-m",
		`{{if eq (getString "format") "emoji"}}{{.gitmoji.Emoji}} {{else}}{{.gitmoji.Code}}{{end}} {{with .scope}}({{.}}): {{end}}{{.title}}`,
		"{{with .message}}-m{{end}}",
		"{{.message}}",
	},
}

var ConventionalCommit = CommitTemplate{
	Questions: []CommitQuestion{
		CommitQuestion{
			PromptType: "conventional",
			Mandatory:  true,
			ValueCode:  "type",
		},
		CommitQuestion{
			PromptType: "text",
			Mandatory:  true,
			Prompt:     "Enter the commit description, with JIRA number at end",
			ValueCode:  "description",
		},
		// TODO: Ask if this is a breaking change
		CommitQuestion{
			PromptType: "text",
			Mandatory:  false,
			Prompt:     "Enter the (optional) commit body",
			ValueCode:  "body",
		},
		CommitQuestion{
			PromptType: "text",
			Mandatory:  false,
			Prompt:     "Enter the (optional) commit footer",
			ValueCode:  "footer",
		},
	},
	Command: "git",
	CommandArgs: []string{
		"commit",
		"-m",
		"{{.type}}: {{.description}}",
		"{{with .body}}-m{{end}}",
		"{{.body}}",
		"{{with .footer}}-m{{end}}",
		"{{.footer}}",
	},
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

		// Disable the pointer (not supported in promptui 0.3.2)
		//Pointer: func(x []rune) []rune { return x },
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func promptGitmoji() (Gitmoji, error) {
	cache, err := NewGitmojiCache()

	if err != nil {
		log.Panic(err)
	}

	gitmoji, err := cache.GetGitmoji()

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
		gitmoji := gitmoji[index]
		tosearch := gitmoji.Name + gitmoji.Code + gitmoji.Description

		// Normalize
		tosearch = strings.Replace(strings.ToLower(tosearch), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(tosearch, input)
	}

	prompt := promptui.Select{
		Label:     "Choose a gitmoji",
		Items:     gitmoji,
		Templates: templates,
		Size:      12,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return Gitmoji{}, err
	}

	return gitmoji[i], nil
}

type ConventionalCommitType struct {
	Name               string
	Description        string
	IncludeInChangelog bool
}

var ConventialCommitTypeList = []ConventionalCommitType{
	ConventionalCommitType{
		Name:               "feat",
		Description:        "A new feature.",
		IncludeInChangelog: true,
	},
	ConventionalCommitType{
		Name:               "fix",
		Description:        "A bug fix.",
		IncludeInChangelog: true,
	},
	ConventionalCommitType{
		Name:               "docs",
		Description:        "Documentation only changes.",
		IncludeInChangelog: false,
	},
	ConventionalCommitType{
		Name:               "perf",
		Description:        "A code change that improves performance.",
		IncludeInChangelog: true,
	},
	ConventionalCommitType{
		Name:               "refactor",
		Description:        "A code change that neither fixes a bug nor adds a feature.",
		IncludeInChangelog: false,
	},
	ConventionalCommitType{
		Name:               "test",
		Description:        "Adding missing or correcting existing tests.",
		IncludeInChangelog: false,
	},
	ConventionalCommitType{
		Name:               "chore",
		Description:        "Changes to the build process or auxiliary tools and libraries such as documentation generation.",
		IncludeInChangelog: false,
	},
}

func promptConventionalCommitType() (ConventionalCommitType, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ \"?\" | yellow }} {{ . }}",
		Active:   "‣ {{ .Name }}",
		Inactive: "  {{ .Name }}",
		Selected: `{{ "? Choose the type of commit" | faint }} {{ .Name }}  - {{ .Description }}`,
		Details: `
--------- Conventional Commit ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Description }}`,
	}

	searcher := func(input string, index int) bool {
		t := ConventialCommitTypeList[index]
		tosearch := t.Name + t.Description

		// Normalize
		tosearch = strings.Replace(strings.ToLower(tosearch), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(tosearch, input)
	}

	prompt := promptui.Select{
		Label:     "Choose the type of commit:",
		Items:     ConventialCommitTypeList,
		Templates: templates,
		Size:      12,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return ConventionalCommitType{}, err
	}

	return ConventialCommitTypeList[i], nil
}
