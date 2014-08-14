package main

import (
	"encoding/json"
	"errors"
	"github.com/izqui/functional"
	"os"
	"os/user"
	"path"
)

const (
	HOME_DIRECTORY_CONFIG = "my home dir"
	TODOS_CONF            = "conf.json"
	ISSUE_CACHE           = "issues.json"
	TOKEN_URL             = "https://github.com/settings/tokens/new?scopes=repo,public_repo"
)

type Configuration struct {
	GithubToken string `json:"github_token,omitempty"`
	Owner       string `json:"github_owner,omitempty"`
	Repo        string `json:"github_repo,omitempty"`
}

type ConfFile struct {
	Config Configuration
	File   *os.File
}

type Issue struct {
	File        string `json:"file,omitempty"`
	Hash        string `json:"hash,omitempty"`
	IssueNumber int    `json:"issue_number,omitempty"`

	Line     int    `json:"-"`
	IssueURL string `json:"-"`
}

type IssueSlice []*Issue

func (slice IssueSlice) Len() int {

	return len(slice)
}

func (slice IssueSlice) Swap(i, j int) {

	slice[i], slice[j] = slice[j], slice[i]
}

func (slice IssueSlice) Less(i, j int) bool {

	return slice[i].Hash < slice[j].Hash
}

func (slice IssueSlice) remove(i int) IssueSlice {

	copy(slice[i:], slice[i+1:])
	slice[len(slice)-1] = nil
	return slice[:len(slice)-1]
}

type IssueCacheFile struct {
	Issues IssueSlice
	File   *os.File
}

func OpenConfiguration(dir string) *ConfFile {

	if dir == HOME_DIRECTORY_CONFIG {

		usr, err := user.Current()
		logOnError(err)
		dir = usr.HomeDir
	}

	dir = path.Join(dir, TODOS_DIRECTORY)

	//Create ~/.todos directory
	err := os.MkdirAll(dir, 0777)
	logOnError(err)

	//Search for conf file inside directory
	f, err := os.OpenFile(path.Join(dir, TODOS_CONF), os.O_RDWR|os.O_CREATE, 0660)
	logOnError(err)

	conf := Configuration{}
	json.NewDecoder(f).Decode(&conf)

	return &ConfFile{File: f, Config: conf}
}

func (conf *ConfFile) WriteConfiguration() error {

	conf.File.Truncate(0)
	conf.File.Seek(0, 0)
	err := json.NewEncoder(conf.File).Encode(conf.Config)
	logOnError(err)

	return conf.File.Close()
}

func IssueCacheFilePath(dir string) string {

	return path.Join(dir, TODOS_DIRECTORY, ISSUE_CACHE)
}
func LoadIssueCache(dir string) *IssueCacheFile {

	filepath := IssueCacheFilePath(dir)

	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0660)
	logOnError(err)

	issues := IssueSlice{}
	json.NewDecoder(f).Decode(&issues)

	return &IssueCacheFile{File: f, Issues: issues}
}

func (cache *IssueCacheFile) GetIssuesInFile(file string) IssueSlice {

	array := IssueSlice{}

	for _, is := range cache.Issues {

		if is != nil && is.File == file {

			array = append(array, is)
		}
	}

	return array

}

func (cache *IssueCacheFile) RemoveIssue(issue Issue) error {

	for i, is := range cache.Issues {

		if is != nil && issue == *is {

			cache.Issues.remove(i)
		}
	}

	return errors.New("Not found")
}

/*
func (cache *IssueCacheFile) GetIssue(file, title string) (Issue, error) {

	for _, i := range cache.Issues {

		if i.File == file && i.Hash ==
		return i, nil
	}

	return Issue{}, errors.New("Not found")
}*/

func (cache *IssueCacheFile) WriteIssueCache() error {

	cache.File.Truncate(0)
	cache.File.Seek(0, 0)

	cache.Issues = functional.Filter(func(i *Issue) bool { return i != nil }, cache.Issues).([]*Issue)
	err := json.NewEncoder(cache.File).Encode(cache.Issues)
	logOnError(err)

	return cache.File.Close()
}
