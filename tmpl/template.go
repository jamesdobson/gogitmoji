package tmpl

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
