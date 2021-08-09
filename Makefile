GIT_COMMIT=$(shell git rev-list -1 HEAD)

worklog:
	@# CGO_ENABLED=0 go build -o bin/worklog -ldflags "-X main.GitCommit=$(GIT_COMMIT) -linkmode external -extldflags -static" github.com/husio/worklog/cmd/worklog
	CGO_ENABLED=0 go build -o bin/worklog -ldflags "-X main.GitCommit=$(GIT_COMMIT)" github.com/husio/worklog/cmd/worklog

install: worklog
	cp bin/worklog $(HOME)/.bin/worklog

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

.PHONY: worklog install
