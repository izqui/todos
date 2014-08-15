# Todos
![Todos](https://github.com/izqui/todos/blob/master/demo.gif)

Get the TODO's in your code directly into Github Issues. Uses pre-commit Git hooks to check for new TODO's every time you commit to the repo.

[Example](#how-to-install-it)

### How to install it? 

At the moment of writing, you need to have Go installed. //TODO: Create distribution binaries [Issue: https://github.com/izqui/todos/issues/41]
```.sh 
go get github.com/izqui/todos
go install github.com/izqui/todos
```
### How to use it?

Just run `todos setup` inside the repo you'd like to track your issues.

Make sure you have `$GOPATH/bin` in your `$PATH`

### How does it work?

Todos installs a git `precommit hook` in your local Git repository, so whenever you are about to commit todos will look for "TODO" tags in comments an submit a Github issue. The issue url is referenced in the code, so you can jump directly there when browsing your code.

In the same way, when you delete the TODO from your code, todos will mark the issue as closed in Github and add a message in your commit description, so you know this was the commit that fixed the issue.

* `todos setup`: 
	* Checks if current directory is a git repository
	* Checks ~/.todos/conf for conf file with github token
	* Asks for github token if it doesn't exist
	* Asks for github owner/repo to know where to post issues
	* Adds `precommit hook` and makes it executable
	* Adds `commit-msg hook` and makes it executable 

* `todos work`: 
	* Checks if current directory is a git repository
	* Gets list of files to check from stdin or git diffs and gets list of changed files (git diff --show-names)
	* Inspects this files looking for "// TODO" 
	* Posts issue to github
	* Saves a local cache in `.todos/issues.json` of the issues it adds.
	* Checks the cache for missing todos and closes issue.
	* Saves a file in `.todos/closed.txt` that the `commit-msg hook` will append to the git message commit file. 
