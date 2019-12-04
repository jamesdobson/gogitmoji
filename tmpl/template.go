package tmpl

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

// CommandTemplate represents a command to execute and user prompts to get
// the command arguments.
type CommandTemplate struct {
	Prompts     []Prompt
	Command     string
	CommandArgs []string
}

// Prompt defines a question to ask the user.
type Prompt struct {
	PromptType    string
	Mandatory     bool
	Prompt        string
	ValueCode     string
	EnableSetting string
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
