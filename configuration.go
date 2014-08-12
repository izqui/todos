package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

const (
	HOME_DIRECTORY_CONFIG = "my home dir"
	TODOS_CONF            = "conf.json"
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

func (conf *ConfFile) WriteConfiguration() {

	err := json.NewEncoder(conf.File).Encode(conf.Config)
	logOnError(err)

	conf.File.Close()
}
