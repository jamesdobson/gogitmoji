package gitmoji

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// GitmojiURL is the address from which to download the list of gitmoji.
const GitmojiURL string = "https://raw.githubusercontent.com/carloscuesta/gitmoji/master/packages/gitmojis/src/gitmojis.json"

// GitmojiDirName is the name of the directory under the user's home directory to store the gitmoji list.
const GitmojiDirName string = ".gitmoji"

// GitmojiFileName is the name of the file to store the list of gitmoji.
const GitmojiFileName string = "gitmojis.json"

// Cache is a local file cache for storing gitmoji.
type Cache struct {
	CacheFile string
	gitmoji   []Gitmoji
	url       string
	download  func(string) ([]byte, error)
}

// gitmojiContainer holds a bunch of Gitmoji, for JSON decoding purposes.
type gitmojiContainer struct {
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

// UpdateCache checks the default URL for new gitmoji and updates the cache
// file in local storage if it is stale.
func UpdateCache() error {
	homedir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("Cannot determine home directory: %v", err)
	}

	cacheFile := path.Join(homedir, GitmojiDirName, GitmojiFileName)

	// Get current list
	currentContent, err := ioutil.ReadFile(cacheFile)

	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("Unable to read local gitmoji cache: %v", err)
		}

		currentContent = []byte{}
	}

	// Fetch latest list
	updatedContent, err := download(GitmojiURL)

	if err != nil {
		return fmt.Errorf("Cannot fetch latest gitmoji: %v", err)
	}

	// Compare
	if !bytes.Equal(currentContent, updatedContent) {
		err = writeCache(cacheFile, updatedContent)

		if err != nil {
			return fmt.Errorf("Unable to write local gitmoji cache: %v", err)
		}

		fmt.Println("List of gitmoji updated! 🎉")
	} else {
		fmt.Println("List of gitmoji is already up to date. 👍")
	}

	return nil
}

// NewCache returns a gitmoji cache using the default URL and local storage.
func NewCache() (Cache, error) {
	homedir, err := os.UserHomeDir()

	if err != nil {
		return Cache{}, fmt.Errorf("Cannot determine home directory: %v", err)
	}

	cacheFile := path.Join(homedir, GitmojiDirName, GitmojiFileName)

	return NewCacheWithURLAndCacheFile(GitmojiURL, cacheFile)
}

// NewCacheWithURLAndCacheFile returns a gitmoji cache of a custom URL and
// local storage location. This method is intended to be used for testing only.
func NewCacheWithURLAndCacheFile(url string, cacheFile string) (Cache, error) {
	return Cache{
		CacheFile: cacheFile,
		url:       url,
		gitmoji:   nil,
		download:  download,
	}, nil
}

func download(url string) ([]byte, error) {
	fmt.Println("🌐  Fetching list of gitmoji...")

	// #nosec G107
	r, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Unable to download gitmoji list (from %s): %v", url, err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to download gitmoji list (from %s): %v", url, r.Status)
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, fmt.Errorf("Unable to download gitmoji list: %v", err)
	}

	return body, nil
}

func writeCache(cacheFile string, content []byte) error {
	err := os.MkdirAll(path.Dir(cacheFile), 0755)

	if err != nil {
		return fmt.Errorf("Unable to create gitmoji cache directory: %v", err)
	}

	err = ioutil.WriteFile(cacheFile, content, 0644)

	if err != nil {
		return fmt.Errorf("Unable to write gitmoji cache: %v", err)
	}

	return nil
}

// GetGitmoji gets the gitmoji list from a local file cache if available;
// otherwise, downloads the latest gitmoji list from github.com.
func (cache *Cache) GetGitmoji() ([]Gitmoji, error) {
	if cache.gitmoji != nil {
		return cache.gitmoji, nil
	}

	content, err := ioutil.ReadFile(cache.CacheFile)

	if err != nil {
		if os.IsNotExist(err) {
			content, err = cache.download(cache.url)

			if err != nil {
				return nil, err
			}

			err = writeCache(cache.CacheFile, content)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Unable to read gitmoji cache: %v", err)
		}
	}

	container := gitmojiContainer{}
	err = json.Unmarshal(content, &container)

	if err != nil {
		return nil, fmt.Errorf("Cannot process gitmoji list; perhaps the file %v is corrupted? Underlying error: %v", cache.CacheFile, err)
	}

	cache.gitmoji = container.Gitmoji

	return cache.gitmoji, nil
}
