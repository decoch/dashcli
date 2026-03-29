.PHONY: test vet lint fmt fmt-check ci

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w $(shell git ls-files '*.go')

fmt-check:
	@output="$$(gofmt -l $(shell git ls-files '*.go'))"; \
	if [ -n "$$output" ]; then \
		echo "Files not formatted with gofmt:"; \
		echo "$$output"; \
		exit 1; \
	fi

ci: fmt-check vet lint test
