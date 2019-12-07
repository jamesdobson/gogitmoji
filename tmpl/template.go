package tmpl

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

// TODO: Clean up structure names and member names (do while writing explanation for them in README)

// CommandTemplate represents a command to execute and user prompts to get
// the command arguments.
type CommandTemplate struct {
	Prompts     []Prompt
	Command     string
	CommandArgs []string
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
