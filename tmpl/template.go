package tmpl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/jamesdobson/gogitmoji/gitmoji"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// CommandTemplate represents a command to execute and user prompts to get
// the command arguments.
type CommandTemplate struct {
	Prompts     []Prompt
	Command     string
	CommandArgs []string
	Messages    []string
}

// Prompt defines a question to ask the user.
type Prompt struct {
	Type      string         `yaml:"Type"`
	Mandatory bool           `yaml:"Mandatory,omitempty"`
	Prompt    string         `yaml:"Prompt,omitempty"`
	Name      string         `yaml:"Name"`
	Condition string         `yaml:"Condition,omitempty"`
	Choices   []PromptChoice `yaml:"Choices,omitempty"`
}

// PromptChoice defines a single option in a multiple-choice prompt.
type PromptChoice struct {
	Value       string `yaml:"Value"`
	Description string `yaml:"Description"`
}

// TemplateLookup maps template names to templates.
var TemplateLookup = make(map[string]CommandTemplate, 2)

// DefaultTemplateName is the name of the template to use if when no template is specified.
var DefaultTemplateName = gitmojiCommandTemplateName

// LoadTemplates reads a map of template names to basic data types and populates
// TemplateLookup with the result.
func LoadTemplates(templates map[string]interface{}) {
	for name, t := range templates {
		var result CommandTemplate

		err := mapstructure.Decode(t, &result)

		if err != nil {
			log.Fatalf("Error processing template '%s': %v", name, err)
		}

		TemplateLookup[name] = result
	}
}

// RunTemplateCommand prompts the user for the template prompts and then runs
// the command specified in the template.
func RunTemplateCommand(tpl CommandTemplate) {
	answers := getAnswers(tpl)
	args := generateArgs(&(tpl.CommandArgs), answers)
	displayCommand := getPrintableCommand(tpl.Command, args)
	fmt.Printf("Going to execute: %s\n\n", displayCommand)

	if !confirm("Execute") {
		fmt.Printf("Canceled.\n")
		return
	}

	run(tpl.Command, args)
}

// GetTemplateMessage prompts the user for the teamplate prompts and then
// returns the result of formatting the "Messages" into a string.
func GetTemplateMessage(tpl CommandTemplate) string {
	answers := getAnswers(tpl)
	messages := generateArgs(&(tpl.Messages), answers)

	return strings.Join(messages, "\n\n")
}

func getAnswers(tpl CommandTemplate) map[string]interface{} {
	var answers = map[string]interface{}{}

	for q := 0; q < len(tpl.Prompts); q++ {
		var question = tpl.Prompts[q]

		if question.Condition != "" && !viper.GetBool(question.Condition) {
			continue
		}

		switch question.Type {
		case "text":
			answer := promptOrCancel(question.Prompt, question.Mandatory)
			answers[question.Name] = answer

		case "choice":
			answer := promptChoice(question)
			answers[question.Name] = answer

		case "gitmoji":
			gitmoji, err := promptGitmoji()

			if err != nil {
				if err == promptui.ErrInterrupt {
					fmt.Println("Canceled.")
					os.Exit(1)
				}

				log.Panic("Couldn't pick a gitmoji: ", err)
			}

			answers[question.Name] = gitmoji

		default:
			log.Fatalf("Unknown prompt type '%s'...\n", question.Type)
		}
	}

	return answers
}

func generateArgs(templates *[]string, answers map[string]interface{}) []string {
	var args = make([]string, 0, len(*templates))
	var sb strings.Builder
	var functions = map[string]interface{}{
		"getString": viper.GetString,
	}

	for n := 0; n < len(*templates); n++ {
		sb.Reset()
		t, err := template.New("arg").
			Funcs(functions).
			Parse((*templates)[n])

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
		//Pointer: func(x []rune) []rune { return x },
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
		log.Fatal("Unable to get list of gitmoji: ", err)
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

func promptChoice(question Prompt) string {
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
		Label:     question.Prompt,
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
