package tmpl

var gitmojiCommandTemplateName = "gitmoji"
var gitmojiCommandTemplate = CommandTemplate{
	Prompts: []Prompt{
		{
			Type:      "gitmoji",
			Mandatory: true,
			Name:      "gitmoji",
		},
		{
			Type:      "text",
			Mandatory: false,
			Prompt:    "Enter the scope of current changes",
			Name:      "Scope",
			Condition: "scope",
		},
		{
			Type:      "text",
			Mandatory: true,
			Prompt:    "Enter the commit title",
			Name:      "title",
		},
		{
			Type:      "text",
			Mandatory: false,
			Prompt:    "Enter the (optional) commit message",
			Name:      "message",
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
	Messages: []string{
		`{{if eq (getString "format") "emoji"}}{{.gitmoji.Emoji}} {{else}}{{.gitmoji.Code}}{{end}} {{with .scope}}({{.}}): {{end}}{{.title}}`,
		"{{.message}}",
	},
}

func init() {
	TemplateLookup[gitmojiCommandTemplateName] = gitmojiCommandTemplate
}
