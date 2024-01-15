set export
init:
   if [ "$$(basename "$$(pwd)")" != "chronio-core" ]; then \
        cd "$$(dirname "$$(realpath "$${BASH_SOURCE[0]}")")"; \
    fi
check:
	go build -o /dev/null ./...
test:
	go test ./... -v -short
