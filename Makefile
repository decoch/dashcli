.PHONY: test lint fmt fmt-check ci

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofumpt -w .

fmt-check:
	gofumpt -l . | grep . && exit 1 || exit 0

ci: fmt-check lint test

