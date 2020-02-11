GIT_COMMIT=$(shell git rev-list -1 HEAD)
BUILD_TIME=$(shell date "+%s")

worklog-parse:
	go build -o bin/worklog-parse -ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -linkmode external -extldflags -static" github.com/husio/worklog

help:
	@echo
	@echo "Commands"
	@echo "========"
	@echo
	@sed -n '/^[a-zA-Z0-9_-]*:/s/:.*//p' < Makefile | grep -v -E 'default|help.*' | sort
	@echo
	@echo "Flags"
	@echo "====="
	@echo
	@echo "tag: deploy with custom tag (default: empty or 'tainted')"

.PHONY: worklog-parse
