# Todos

Get the todos in your code directly in Github Issues.

### How does it work?

Todos installs a git precommit hook in your local Git repository, so whenever you are about to commit todos will look for "TODO" tags in comments an submit a Github issue. The issue url is referenced in the code, so you can jump directly there when browsing your code.

In the same way, when you delete the TODO from your code, todos will mark the issue as closed in Github.

* `todos setup`: 
	* checks if dir is a git repo
	* checks ~/.todos/conf for conf file with github token
	* asks for github token if it doesn't exist
	* asks for github owner/repo to know where to post issues
	* adds precommit hook and makes it executable

* `todos work`: 
	* cd's root directory (git rev-parse --show-toplevel) and fail if there's no git repo
	* gets list of files to check from stdin or git diffs and gets list of changed files (git diff --show-names)
	* inspects this files looking for "// TODO" 
	* posts issue to github


## .todos directory
	conf.json ->
