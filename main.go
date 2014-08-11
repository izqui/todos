package main

import (
	"flag"
	"fmt"
	"github.com/izqui/functional"
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

		if len(flag.Args()) < 1 {
			showHelp()
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

				// Set git "precommit" Hook

			case "work":
				diff, _ := GitDiffFiles()

				diff = functional.Map(func(s string) string { return root + "/" + s }, diff).([]string)

				fmt.Println("Files to check: ", diff)

			default:
				showHelp()
			}
		}
	}
}

func showHelp() {

	fmt.Println("Unknown command") //TODO: Write help
}
func logOnError(err error) {

	if err != nil {
		log.Fatal(err)
	}
}
