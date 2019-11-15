package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
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
	gitmoji := getGitmoji()

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
		gitmoji := gitmoji[index]
		tosearch := gitmoji.Name + gitmoji.Code + gitmoji.Description

		// Normalize
		tosearch = strings.Replace(strings.ToLower(tosearch), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(tosearch, input)
	}

	prompt := promptui.Select{
		Label:     "Gitmoji",
		Items:     gitmoji,
		Templates: templates,
		Size:      12,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("%s\n", gitmoji[i].Name)
}

const GitmojiURL string = `https://raw.githubusercontent.com/carloscuesta/gitmoji/master/src/data/gitmojis.json`
const GitmojiDirName string = `.gitmoji`
const GitmojiFileName string = `gitmojis.json`

// Gets the gitmoji list from a local file cache if available;
// otherwise, downloads the latest gitmoji list from github.com.
func getGitmoji() []Gitmoji {
	homedir, err := os.UserHomeDir()

	if err != nil {
		log.Panic(`Cannot determine home directory: `, err)
	}

	cacheFile := path.Join(homedir, GitmojiDirName, GitmojiFileName)

	content, err := ioutil.ReadFile(cacheFile)

	if err != nil {
		if os.IsNotExist(err) {
			log.Print(`Fetching list of gitmoji...`)

			r, err := http.Get(GitmojiURL)

			if err != nil {
				log.Panic(`Unable to download gitmoji list: `, err)
			}

			defer r.Body.Close()
			content, err = ioutil.ReadAll(r.Body)

			if err != nil {
				log.Panic(`Unable to download gitmoji list: `, err)
			}

			err = os.MkdirAll(path.Dir(cacheFile), 0755)

			if err != nil {
				log.Panic(`Unable to create gitmoji cache directory: `, err)
			}

			err = ioutil.WriteFile(cacheFile, content, 0644)

			if err != nil {
				log.Panic(`Unable to write gitmoji cache: `, err)
			}
		} else {
			log.Panic(`Unable to read gitmoji cache: `, err)
		}
	}

	gitmojiContainer := GitmojiContainer{}
	err = json.Unmarshal([]byte(content), &gitmojiContainer)

	if err != nil {
		log.Panic(`Cannot process gitmoji list; perhaps the file is corrupted? `, err)
	}

	return gitmojiContainer.Gitmoji
}
