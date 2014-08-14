# Todos

* `todos setup`: 
	* checks if dir is a git repo
	* checks ~/.todos/conf for conf file with github token
	* asks for github token if it doesn't exist
	* asks for github owner/repo to know where to post issues
	* adds precommit hook and makes it executable

* `todos work`: 
	* cd's root directory (git rev-parse --show-toplevel) and fail if there's no git repo
	* gets list of files to check from stdin or git diffs and gets list of changed files (git diff --show-names)
	* inspects this files looking for "// TODO:" [Issue: https://github.com/izqui/todos/issues/28]
	* posts issue to github


## .todos directory
	conf.json ->
