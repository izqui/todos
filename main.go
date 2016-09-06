package main

import (
	"flag"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"

	"github.com/izqui/functional"
	"github.com/izqui/helpers"
)

var (
	tokenArg = flag.String("token", "", "Github setup token")
	resetArg = flag.Bool("reset", false, "Reset Github token")
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
			case "setup":
				setup(root)

			case "work":

				// Try to read lines from Stdin, if not talk to git
				diff, err := ReadStdin()
				if len(diff) == 0 {

					diff, err = GitDiffFiles()
					logOnError(err)
				}

				diff = functional.Map(func(s string) string { return path.Join(root, s) }, diff).([]string)

				work(root, diff)

			default:
				showHelp()
			}
		}
	}
}

func setup(root string) {

	// Global configuration
	global := OpenConfiguration(HOME_DIRECTORY_CONFIG)
	defer global.File.Close()

	if *tokenArg != "" {

		global.Config.GithubToken = *tokenArg

	} else if global.Config.GithubToken == "" || *resetArg {

		fmt.Println("Paste Github access token:")
		open.Run(TOKEN_URL)
		var scanToken string
		fmt.Scanln(&scanToken)
		global.Config.GithubToken = scanToken //TODO: Check if token is valid [Issue: https://github.com/izqui/todos/issues/32]
	}
	logOnError(global.WriteConfiguration())

	// Local configuration
	local := OpenConfiguration(root)
	defer local.File.Close()

	if local.Config.Owner == "" || local.Config.Repo == "" || *resetArg {

		owner, repo, _ := GitGetOwnerRepoFromRepository()
		fmt.Printf("Enter the Github owner of the repo (Default: %s):\n", owner)
		fmt.Scanln(&owner)
		fmt.Printf("Enter the Github repo name (Default: %s):\n", repo)
		fmt.Scanln(&repo)

		// TODO: Check if repository exists [Issue: https://github.com/izqui/todos/issues/33]
		local.Config.Owner = owner
		local.Config.Repo = repo
	}

	logOnError(local.WriteConfiguration())
	logOnError(GitAdd(path.Join(root, TODOS_DIRECTORY)))

	SetupGitPrecommitHook(root)
	SetupGitCommitMsgHook(root)
}

func work(root string, files []string) {

	//Load configuration
	global := OpenConfiguration(HOME_DIRECTORY_CONFIG)
	defer global.File.Close()
	local := OpenConfiguration(root)
	defer local.File.Close()

	if global.Config.GithubToken == "" {

		fmt.Println("[Todos] Missing Github token. Set it in ~/.todos/conf.json")

	} else if local.Config.Owner == "" || local.Config.Repo == "" {

		fmt.Println("[Todos] You need to setup the repo running 'todos setup'")

	} else {
		
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: global.Config.GithubToken},
		)
		tc := oauth2.NewClient(oauth2.NoContext, ts)

		owner := local.Config.Owner
		repo := local.Config.Repo

		fmt.Printf("[Todos] Scanning changed files on %s/%s\n", owner, repo)

		client := github.NewClient(tc)

		existingRegex, err := regexp.Compile(ISSUE_URL_REGEX)
		logOnError(err)
		todoRegex, err := regexp.Compile(TODO_REGEX)
		logOnError(err)

		cacheFile := LoadIssueCache(root)
		cacheChanges := false

		//Leave first two lines blank so it displays as the commit description
		closedIssues := []string{"", ""}

		for _, file := range files {

			// In the cache files are saved as a relative path to the project root
			relativeFilePath := pathDifference(root, file)

			fileIssuesCache := cacheFile.GetIssuesInFile(relativeFilePath)
			fileIssuesCacheCopy := fileIssuesCache

			removed := 0

			fmt.Println("[Todos] Checking file:", relativeFilePath)

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
							fileIssuesCacheCopy = fileIssuesCacheCopy.remove(i)
							removed++
						}
					}

				} else if todo != "" {

					issuesCount++
					go func(line int, cb chan Issue) {

						branch, _ := GitBranch()

						filename := path.Base(file)

						body := fmt.Sprintf(ISSUE_BODY, filename, fmt.Sprintf(GITHUB_FILE_URL, owner, repo, branch, relativeFilePath))
						issue, _, err := client.Issues.Create(owner, repo, &github.IssueRequest{Title: &todo, Body: &body})
						logOnError(err)

						if issue != nil {
							cb <- Issue{IssueURL: *issue.HTMLURL, IssueNumber: *issue.Number, Line: line, File: relativeFilePath}
						}
					}(i, cb)
				}
			}
		loop:
			for issuesCount > 0 {

				select {
				case issue := <-cb:

					ref := fmt.Sprintf("[Issue: %s]", issue.IssueURL)
					lines[issue.Line] = fmt.Sprintf("%s %s", lines[issue.Line], ref)
					fmt.Printf("[Todos] Created issue #%d\n", issue.IssueNumber)
					changes = true
					issuesCount--

					issue.Hash = helpers.SHA1([]byte(ref))
					cacheFile.Issues = append(cacheFile.Issues, &issue)
					cacheChanges = true

				case <-timeout(3 * time.Second):
					break loop
				}
			}

			closeCount := 0
			closeCb := make(chan Issue)
			for _, is := range fileIssuesCacheCopy {

				if is != nil {

					closeCount++
					go func(i Issue) {

						var closed string = "closed"
						_, _, err := client.Issues.Edit(owner, repo, is.IssueNumber, &github.IssueRequest{State: &closed})
						logOnError(err)
						closeCb <- i
					}(*is)
				}
			}

		loops:
			for closeCount > 0 {
				select {
				case is := <-closeCb:
					closeCount--
					issueClosing := fmt.Sprintf("[Todos] Closed issue #%d", is.IssueNumber)
					fmt.Println(issueClosing)
					closedIssues = append(closedIssues, issueClosing)
					cacheFile.RemoveIssue(is)
					cacheChanges = true

				case <-timeout(3 * time.Second):
					break loops
				}
			}

			if changes {
				logOnError(WriteFileLines(file, lines, false))
				GitAdd(file)
			} else {
				fmt.Println("[Todos] No new todos found")
			}
		}

		if cacheChanges {
			logOnError(cacheFile.WriteIssueCache())
			GitAdd(IssueCacheFilePath(root))
		}
		if len(closedIssues) <= 2 {

			closedIssues = []string{}
		}

		logOnError(WriteFileLines(path.Join(root, TODOS_DIRECTORY, CLOSED_ISSUES_FILENAME), closedIssues, false))
	}
}

func pathDifference(p1, p2 string) string {

	return path.Join(strings.Split(p2, "/")[len(strings.Split(p1, "/")):]...)
}

func timeout(i time.Duration) chan bool {

	t := make(chan bool)
	go func() {
		time.Sleep(i)
		t <- true
	}()

	return t
}

func showHelp() {

	fmt.Println("Unknown command.")
	fmt.Println("\t* setup: Setup the current repository.")
	fmt.Println("\t* work: Runs todos and looks for todos in files in current git diff.")
}
func logOnError(err error) {

	if err != nil {
		log.Println("[Todos] Err:", err)
	}
}
