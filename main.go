package main // Inserted comment

import (
	"flag"
	"fmt"
	"github.com/izqui/functional"
	"github.com/skratchdot/open-golang/open"
	"log"
	"regexp"
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

				fmt.Println("Scanning changed files...")
				diff, _ := GitDiffFiles()
				diff = functional.Map(func(s string) string { return root + "/" + s }, diff).([]string)

				existingRegex, err := regexp.Compile(ISSUE_URL_REGEX)
				logOnError(err)
				todoRegex, err := regexp.Compile(TODO_REGEX)
				logOnError(err)

				for _, file := range diff {
					lines, err := ReadFileLines(file)
					logOnError(err)

					for _, line := range lines {
						ex := existingRegex.FindString(line)
						todo := todoRegex.FindString(line)

						if ex == "" && todo != "" {
							fmt.Println("Found todo", todo)
						}
					}
				}

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
