package main

import (
	"flag"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"log"
)

func init() {

	flag.Parse()
}

func main() {

	root, err := GitDirectoryRoot()

	if err != nil {

		fmt.Println("You must use todos inside a git repository")

	} else {

		mode := flag.Args()[0]
		switch mode {
		case "install":

			f := OpenConfiguration(HOME_DIRECTORY_CONFIG)
			defer f.File.Close()

			// Check config file for github access token
			if f.Config.GithubToken == "" {

				fmt.Println("Paste Github access token:")
				open.Run(TOKEN_URL)
				var token string
				fmt.Scanln(&token)
				f.Config.GithubToken = token //TODO: Check if token is valid
				f.WriteConfiguration()
			}

			// Set Git Hook

		case "work":
			diff, err := GitDiffFiles()
			fmt.Println(diff)

		default:
			fmt.Println("Unknown command")
		}
	}
}

func logOnError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
