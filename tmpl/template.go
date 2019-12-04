package tmpl

import (
	"log"

	"github.com/mitchellh/mapstructure"
)

// TODO: Add comments to fix linter warnings
type CommitTemplate struct {
	Questions   []CommitQuestion
	Command     string
	CommandArgs []string
}

type CommitQuestion struct {
	PromptType    string
	Mandatory     bool
	Prompt        string
	ValueCode     string
	EnableSetting string
}

var TemplateLookup = make(map[string]CommitTemplate, 2)
var DefaultTemplateName = gitmojiCommitTemplateName

func LoadTemplates(templates map[string]interface{}) {
	for name, t := range templates {
		var result CommitTemplate

		err := mapstructure.Decode(t, &result)

		if err != nil {
			log.Fatalf("Error processing template '%s': %v", name, err)
		}

		TemplateLookup[name] = result
	}
}
