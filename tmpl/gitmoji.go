package tmpl

var gitmojiCommitTemplateName = "gitmoji"
var gitmojiCommitTemplate = CommitTemplate{
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
			EnableSetting: "scope",
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

func init() {
	TemplateLookup[gitmojiCommitTemplateName] = gitmojiCommitTemplate
}
