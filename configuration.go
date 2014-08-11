package main

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
)

const (
	HOME_DIRECTORY_CONFIG = "my home dir"
	TODOS_DIRECTORY       = ".todos"
	TODOS_CONF            = "conf.json"
	TOKEN_URL             = "https://github.com/settings/tokens/new?scopes=repo,public_repo"
)

type Configuration struct {
	GithubToken string `json:"github_token"`
}

type ConfFile struct {
	Config Configuration
	File   *os.File
}

func OpenConfiguration(path string) *ConfFile {

	if path == HOME_DIRECTORY_CONFIG {

		//Create ~/.todos directory
		usr, err := user.Current()
		logOnError(err)
		todosDir := usr.HomeDir + "/.todos"

		err = os.MkdirAll(todosDir, 0777)
		logOnError(err)

		//Search for conf file inside directory
		f, err := os.OpenFile(todosDir+"/conf.json", os.O_RDWR|os.O_CREATE, 0660)
		logOnError(err)
		//defer f.Close()

		conf := Configuration{}
		json.NewDecoder(f).Decode(&conf)

		return &ConfFile{File: f, Config: conf}
	} else {

		log.Fatal("Only supporting config in home directory ATM")
	}

	return nil
}

func (conf *ConfFile) WriteConfiguration() {

	err := json.NewEncoder(conf.File).Encode(conf.Config)
	logOnError(err)

	conf.File.Close()
}
