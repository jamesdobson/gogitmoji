package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/manifoldco/promptui"
)

type GitmojiContainer struct {
	Gitmoji []Gitmoji `json:"gitmojis"`
}

type Gitmoji struct {
	Emoji       string
	Entity      string
	Code        string
	Description string
	Name        string
}

func main() {
	file, _ := ioutil.ReadFile("gitmojis.json")
	gitmojiContainer := GitmojiContainer{}
	_ = json.Unmarshal([]byte(file), &gitmojiContainer)

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "‣ {{ .Emoji }} {{ .Code | cyan }} {{ .Description }}",
		Inactive: "  {{ .Emoji }} {{ .Code | cyan }} {{ .Description }}",
		Selected: "{{ .Emoji }} {{ .Code | cyan }} {{ .Description }}",
		Details: `
--------- Gitmoji ----------
{{ "Name:" | faint }}	{{ .Emoji }} {{ .Name }}
{{ "Entity:" | faint }}	{{ .Entity }}
{{ "Code:" | faint }}	{{ .Code }}
{{ "Description:" | faint }}	{{ .Description }}`,
	}

	searcher := func(input string, index int) bool {
		gitmoji := gitmojiContainer.Gitmoji[index]
		tosearch := gitmoji.Name + gitmoji.Code + gitmoji.Description

		// Normalize
		tosearch = strings.Replace(strings.ToLower(tosearch), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(tosearch, input)
	}

	prompt := promptui.Select{
		Label:     "Gitmoji",
		Items:     gitmojiContainer.Gitmoji,
		Templates: templates,
		Size:      12,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("%s\n", gitmojiContainer.Gitmoji[i].Name)
}
