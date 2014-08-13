package main

//TODO: Todo creation is working again :) [Issue: https://github.com/izqui/todos/issues/9]
import (
	"flag"
	"fmt"
	"log"
	"path"
	"regexp"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
	"github.com/skratchdot/open-golang/open"

	"github.com/izqui/functional"
	"github.com/izqui/helpers"
)

var (
	token = flag.String("token", "", "Github setup token")
	reset = flag.Bool("reset", false, "Reset Github token")
)

const (
	TODOS_DIRECTORY = ".todos/"
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

			global := OpenConfiguration(HOME_DIRECTORY_CONFIG)
			defer global.File.Close()
			local := OpenConfiguration(root)
			defer local.File.Close()

			mode := flag.Args()[0]
			switch mode {
			case "setup":

				// Check config file for github access token
				if *token != "" {
					global.Config.GithubToken = *token

				} else if global.Config.GithubToken == "" || *reset {

					fmt.Println("Paste Github access token:")
					open.Run(TOKEN_URL)
					var scanToken string
					fmt.Scanln(&scanToken)
					global.Config.GithubToken = scanToken //TODO: Check if token is valid [Issue: https://github.com/izqui/todos/issues/5]
				}
				global.WriteConfiguration()

				if local.Config.Owner == "" || local.Config.Repo == "" || *reset {

					owner, repo, _ := GitGetOwnerRepoFromRepository()
					fmt.Printf("Enter the Github owner of the repo (Default: %s):\n", owner)
					fmt.Scanln(&owner)
					fmt.Printf("Enter the Github repo name (Default: %s):\n", repo)
					fmt.Scanln(&repo)

					// TODO: Check if repository exists [Issue: https://github.com/izqui/todos/issues/8]
					local.Config.Owner = owner
					local.Config.Repo = repo
				}

				logOnError(local.WriteConfiguration())
				logOnError(GitAdd(path.Join(root, TODOS_DIRECTORY)))

				setupHook(root + "/.git/hooks/pre-commit")

			case "work":

				if global.Config.GithubToken == "" {

					fmt.Println("Missing Github token. Set it in ~/.todos/conf.json")

				} else {

					o := &oauth.Transport{
						Token: &oauth.Token{AccessToken: global.Config.GithubToken},
					}

					owner := local.Config.Owner
					repo := local.Config.Repo

					fmt.Printf("[Todos] Scanning changed files on %s/%s\n", owner, repo)

					client := github.NewClient(o.Client())

					// Try to read lines from Stdin, if not talk to git
					diff, err := ReadStdin()
					if len(diff) == 0 {

						diff, err = GitDiffFiles()
						logOnError(err)
					}

					diff = functional.Map(func(s string) string { return path.Join(root, s) }, diff).([]string)

					existingRegex, err := regexp.Compile(ISSUE_URL_REGEX)
					logOnError(err)
					todoRegex, err := regexp.Compile(TODO_REGEX)
					logOnError(err)

					cacheFile := LoadIssueCache(root)
					cacheChanges := false

					for _, file := range diff {

						fileIssuesCache := cacheFile.GetIssuesInFile(file)
						fileIssuesCacheCopy := fileIssuesCache
						removed := 0

						fmt.Println("[Todos] Checking file: ", file)

						lines, err := ReadFileLines(file)
						logOnError(err)

						changes := false

						cb := make(chan Issue)
						issuesCount := 0

						for i, line := range lines {

							ex := existingRegex.FindString(line)
							todo := todoRegex.FindString(line)

							if ex != "" {

								for i, is := range fileIssuesCache {

									if is != nil && is.Hash == helpers.SHA1([]byte(ex)) {

										cacheChanges = true
										fileIssuesCacheCopy.remove(i - removed)
										removed++
										fmt.Println("Detected Existing Issue #", is.IssueNumber)
									}
								}

							} else if todo != "" {

								issuesCount++
								go func(line int, cb chan Issue) {

									issue, _, err := client.Issues.Create(owner, repo, &github.IssueRequest{Title: &todo})
									logOnError(err)
									cb <- Issue{IssueURL: *issue.HTMLURL, IssueNumber: *issue.Number, Line: line, File: file}
								}(i, cb)
							}
						}
					loop:
						for issuesCount > 0 {

							select {
							case issue := <-cb:

								ref := fmt.Sprintf("[Issue: %s]", issue.IssueURL)
								lines[issue.Line] = fmt.Sprintf("%s %s", lines[issue.Line], ref)
								fmt.Printf("[Todos] Created issue #%d: %s\n", issue.IssueNumber)
								changes = true
								issuesCount--

								issue.Hash = helpers.SHA1([]byte(ref))
								cacheFile.Issues = append(cacheFile.Issues, &issue)
								cacheChanges = true

							case <-timeout(3 * time.Second):
								break loop
							}
						}

						for _, is := range fileIssuesCacheCopy {

							if is != nil {

								fmt.Println("Detected removed Issue #", is.IssueNumber)
							}
						}

						if changes {
							logOnError(WriteFileLines(file, lines, false))
						} else {
							fmt.Println("[Todos] No todos found")
						}

						if cacheChanges {
							logOnError(cacheFile.WriteIssueCache())
						}
					}
				}

			default:
				showHelp()
			}
		}
	}
}

func timeout(i time.Duration) chan bool {

	t := make(chan bool)
	go func() {
		time.Sleep(i)
		t <- true
	}()

	return t
}

func setupHook(path string) {

	bash := "#!/bin/bash"
	script := "git diff --cached --name-only | todos work"

	lns, err := ReadFileLines(path)
	logOnError(err)

	if len(lns) == 0 {
		lns = append(lns, bash)
	}

	//Filter existing script line
	lns = functional.Filter(func(a string) bool { return a != script }, lns).([]string)
	lns = append(lns, script)

	logOnError(WriteFileLines(path, lns, true))

}

func showHelp() {

	fmt.Println("Unknown command") //TODO: Write help [Issue: https://github.com/izqui/todos/issues/7]
}
func logOnError(err error) {

	if err != nil {
		log.Println("Err:", err)
	}
}
