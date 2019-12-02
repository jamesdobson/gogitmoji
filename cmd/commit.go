package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	t, err := template.New("test").Parse(GitmojiCommit.Messages[0])

	if err != nil {
		panic(err)
	}

	var d = map[string]string{
		"gitmoji": "X",
		"title":   "this is a test title",
		//"scope":   "mini",
	}

	err = t.Execute(os.Stdout, d)
}

/*
func commit() {
	t := viper.GetString(typeSetting)

	if t == typeGitmoji {
		commitAsGitmoji()
	} else if t == typeConventional {
		commitAsConventional()
	} else {
		fmt.Printf("Unknown commit type: \"%s\"\n", t)
		os.Exit(1)
	}

	fmt.Printf("\ngogitmoji done.\n")
}
*/

func execGit(args []string) {
	fmt.Printf("Executing...\n")

	var cmd = exec.Command("git", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		code := cmd.ProcessState.ExitCode()
		fmt.Printf("\ngit commit exited with code: %d\n", code)
		os.Exit(code)
	}
}

type CommitQuestion struct {
	PromptType string
	Mandatory  bool
	Prompt     string
	ValueCode  string
}

type CommitTemplate struct {
	Questions []CommitQuestion
	Messages  []string
}

var GitmojiCommit = CommitTemplate{
	Questions: []CommitQuestion{
		CommitQuestion{
			PromptType: "gitmoji",
			Mandatory:  true,
			ValueCode:  "gitmoji",
		},
		// TODO: How to make this driven by config?
		CommitQuestion{
			PromptType: "text",
			Mandatory:  false,
			Prompt:     "Enter the scope of current changes",
			ValueCode:  "Scope",
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
	Messages: []string{
		"{{ .gitmoji }}  {{ with .scope }}({{ . }}): {{ end }}{{ .title }}\n\n",
		"{{ message }}",
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
	Messages: []string{
		"{{ type }}: {{ description }}",
		"{{ body }}",
		"{{ footer }}",
	},
}

func commitAsGitmoji() {
	format := viper.GetString(formatSetting)

	if format != formatAsEmoji && format != formatAsCode {
		fmt.Printf("Unknown emoji format: \"%s\"\n", format)
		os.Exit(1)
	}

	cache, err := NewGitmojiCache()

	if err != nil {
		log.Panic(err)
	}

	gitmojiList, err := cache.GetGitmoji()

	if err != nil {
		log.Panic("Unable to get list of gitmoji: ", err)
	}

	gitmoji, err := promptGitmoji(gitmojiList)

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Canceled.")
			os.Exit(1)
		}

		log.Panic("Couldn't pick a gitmoji: ", err)
	}

	var scope string = ""

	if viper.GetBool(scopeSetting) {
		scope = promptOrCancel("Enter the scope of current changes", true)
	}

	title := promptOrCancel("Enter the commit title", true)
	message := promptOrCancel("Enter the (optional) commit message", false)

	var fullTitle string = title

	if scope != "" {
		fullTitle = "(" + scope + "): " + title
	}

	if format == formatAsEmoji {
		fullTitle = gitmoji.Emoji + "  " + fullTitle
	} else {
		fullTitle = gitmoji.Code + "  " + fullTitle
	}

	fmt.Printf("Going to execute:\n\ngit commit -m \"%s\"", fullTitle)

	var args []string = make([]string, 3, 5)

	args[0] = "commit"
	args[1] = "-m"
	args[2] = fullTitle

	if len(message) > 0 {
		fmt.Printf(" -m \"%s\"", message)
		args = append(args, "-m", message)
	}

	fmt.Printf("\n\n")

	if !confirm("Execute") {
		fmt.Printf("Canceled.\n")
		return
	}

	execGit(args)
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

func promptGitmoji(gitmoji []Gitmoji) (Gitmoji, error) {
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

/*
Commit message pattern

<type>: <description> <jira number>

[optional body]

[optional footer]

Footer must start with "BREAKING CHANGE: " if there is a breaking change, and a
description must follow.
*/

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

func commitAsConventional() {
	t, err := promptConventionalCommitType()

	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Canceled.")
			os.Exit(1)
		}

		log.Panic("Couldn't pick a conventional commit type: ", err)
	}

	description := promptOrCancel("Enter the commit description, with JIRA number at end", true)
	body := promptOrCancel("Enter the commit body", false)
	footer := promptOrCancel("Enter the commit footer", false)

	var fullDescription string = t.Name + ": " + description

	fmt.Printf("Going to execute:\n\ngit commit -m \"%s\"", fullDescription)

	var args []string = make([]string, 3, 7)

	args[0] = "commit"
	args[1] = "-m"
	args[2] = fullDescription

	if len(body) > 0 {
		fmt.Printf(" -m \"%s\"", body)
		args = append(args, "-m", body)
	}

	if len(footer) > 0 {
		fmt.Printf(" -m \"%s\"", footer)
		args = append(args, "-m", footer)
	}

	fmt.Printf("\n\n")

	if !confirm("Execute") {
		fmt.Printf("Canceled.\n")
		return
	}

	execGit(args)
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
