package main

const (
	ISSUE_URL_REGEX = "\\[(Issue:[^\\]]*)\\]"
	TODO_REGEX      = "TODO(\\([^)]+\\))?:(.*)"
	ISSUE_BODY      = "On file: [%s](%s)"
	GITHUB_FILE_URL = "https://github.com/%s/%s/blob/%s/%s"
)
