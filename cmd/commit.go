package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	formatSetting = "format"

	formatAsEmoji = "emoji"
	formatAsCode  = "code"

	scopeSetting = "scope"
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

	commitCmd.Flags().StringP("format", "f", "emoji", `Emoji format; either "emoji" or "code".`)
	commitCmd.Flags().BoolP("scope", "p", false, "Enable scope prompt")
	viper.BindPFlag(formatSetting, commitCmd.Flags().Lookup("format"))
	viper.BindPFlag(scopeSetting, commitCmd.Flags().Lookup("scope"))
}

func commit() {
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

	if len(message) > 0 {
		fmt.Printf(" -m \"%s\"", message)
	}

	fmt.Printf("\n\n")

	if !confirm("Execute") {
		fmt.Printf("Canceled.\n")
		return
	}

	fmt.Printf("Executing...\n")

	var cmd *exec.Cmd

	if len(message) > 0 {
		cmd = exec.Command("git", "commit", "-m", fullTitle, "-m", message)
	} else {
		cmd = exec.Command("git", "commit", "-m", fullTitle)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		code := cmd.ProcessState.ExitCode()
		fmt.Printf("\ngit commit exited with code: %d\n", code)
		os.Exit(code)
	}

	fmt.Printf("\ngogitmoji done.\n")
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
