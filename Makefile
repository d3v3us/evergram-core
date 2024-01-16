.DEFAULT_GOAL := init
.PHONY: init check test

init:
	@if [ "$$(basename "$$(pwd)")" != "evergram-core" ]; then \
		cd "$$(dirname "$$(realpath "$${BASH_SOURCE[0]}")")"; \
	fi

check:
	go build -o /dev/null ./...

test:
	go test ./... -v -short
