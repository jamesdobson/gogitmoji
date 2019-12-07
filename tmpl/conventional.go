package tmpl

var conventionalCommandTemplateName = "conventional"
var conventionalCommandTemplate = CommandTemplate{
	Prompts: []Prompt{
		{
			PromptType: "choice",
			Prompt:     "Choose the type of commit:",
			Mandatory:  true,
			ValueCode:  "type",
			Choices: []PromptChoice{
				{
					Value:       "feat",
					Description: "A new feature.",
				},
				{
					Value:       "fix",
					Description: "A bug fix.",
				},
				{
					Value:       "docs",
					Description: "Documentation only changes.",
				},
				{
					Value:       "perf",
					Description: "A code change that improves performance.",
				},
				{
					Value:       "refactor",
					Description: "A code change that neither fixes a bug nor adds a feature.",
				},
				{
					Value:       "test",
					Description: "Adding missing or correcting existing tests.",
				},
				{
					Value:       "chore",
					Description: "Changes to the build process or auxiliary tools and libraries such as documentation generation.",
				},
			},
		},
		{
			PromptType: "text",
			Mandatory:  true,
			Prompt:     "Enter the commit description, with JIRA number at end",
			ValueCode:  "description",
		},
		// TODO: Ask if this is a breaking change
		{
			PromptType: "text",
			Mandatory:  false,
			Prompt:     "Enter the (optional) commit body",
			ValueCode:  "body",
		},
		{
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

func init() {
	TemplateLookup[conventionalCommandTemplateName] = conventionalCommandTemplate
}
