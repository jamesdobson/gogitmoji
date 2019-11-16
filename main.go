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

// GitmojiContainer holds a bunch of Gitmoji.
type GitmojiContainer struct {
	Gitmoji []Gitmoji `json:"gitmojis"`
}

// Gitmoji is a structure with the information about a single gitmoji.
type Gitmoji struct {
	Emoji       string
	Entity      string
	Code        string
	Description string
	Name        string
}

func main() {
	cache, err := NewGitmojiCache()

	if err != nil {
		log.Panic(err)
	}

	gitmoji, err := cache.GetGitmoji()

	if err != nil {
		log.Panic("Unable to get list of gitmoji: ", err)
	}

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
		log.Panic("Prompt failed: ", err)
	}

	fmt.Printf("%s\n", gitmoji[i].Name)
}

// GitmojiURL is the address from which to download the list of gitmoji.
const GitmojiURL string = "https://raw.githubusercontent.com/carloscuesta/gitmoji/master/src/data/gitmojis.json"

// GitmojiDirName is the name of the directory under the user's home directory to store the gitmoji list.
const GitmojiDirName string = ".gitmoji"

// GitmojiFileName is the name of the file to store the list of gitmoji.
const GitmojiFileName string = "gitmojis.json"

// GitmojiCache is a local file cache for storing gitmoji.
type GitmojiCache struct {
	CacheFile string
	gitmoji   []Gitmoji
	url       string
	load      func() ([]byte, error)
}

func NewGitmojiCache() (GitmojiCache, error) {
	homedir, err := os.UserHomeDir()

	if err != nil {
		return GitmojiCache{}, fmt.Errorf("Cannot determine home directory: %v", err)
	}

	cacheFile := path.Join(homedir, GitmojiDirName, GitmojiFileName)

	return NewGitmojiCacheWithURLAndCacheFile(GitmojiURL, cacheFile)
}

func NewGitmojiCacheWithURLAndCacheFile(url string, cacheFile string) (GitmojiCache, error) {
	return GitmojiCache{
		CacheFile: cacheFile,
		url:       url,
		gitmoji:   nil,
		load: func() ([]byte, error) {
			r, err := http.Get(url)

			if err != nil {
				return nil, fmt.Errorf("Unable to download gitmoji list: %v", err)
			}

			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				return nil, fmt.Errorf("Unable to download gitmoji list: %v", err)
			}

			return body, nil
		},
	}, nil
}

// GetGitmoji gets the gitmoji list from a local file cache if available;
// otherwise, downloads the latest gitmoji list from github.com.
func (cache *GitmojiCache) GetGitmoji() ([]Gitmoji, error) {
	if cache.gitmoji != nil {
		return cache.gitmoji, nil
	}

	content, err := ioutil.ReadFile(cache.CacheFile)

	if err != nil {
		if os.IsNotExist(err) {
			log.Print("Fetching list of gitmoji...")

			content, err = cache.load()

			if err != nil {
				return nil, err
			}

			err = os.MkdirAll(path.Dir(cache.CacheFile), 0755)

			if err != nil {
				return nil, fmt.Errorf("Unable to create gitmoji cache directory: %v", err)
			}

			err = ioutil.WriteFile(cache.CacheFile, content, 0644)

			if err != nil {
				return nil, fmt.Errorf("Unable to write gitmoji cache: %v", err)
			}
		} else {
			return nil, fmt.Errorf("Unable to read gitmoji cache: %v", err)
		}
	}

	container := GitmojiContainer{}
	err = json.Unmarshal([]byte(content), &container)

	if err != nil {
		return nil, fmt.Errorf("Cannot process gitmoji list; perhaps the file is corrupted? Underlying error: %v", err)
	}

	cache.gitmoji = container.Gitmoji

	return cache.gitmoji, nil
}
