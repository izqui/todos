# Todos

* `todos install`: 
	* checks if it is a git directory 
	* checks ~/.todos/conf for conf file with github token
	* asks for github token if it doesn't exist
	* adds precommit hook and makes it executable

* `todos work`: 
	* cd's root directory (git rev-parse --show-toplevel) and fail if there's no git repo
	* git diffs and gets list of changed files (git diff --show-names)
	* inspects this files looking for "//TODOS:"


## .todos directory
	conf.json ->