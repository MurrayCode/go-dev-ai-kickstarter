.PHONY: fmt test check run

fmt:
	go fmt ./...

test:
	go test ./...

check: fmt test

run:
	go run ./cmd/app
