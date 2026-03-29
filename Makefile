.PHONY: test vet lint fmt fmt-check ci

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

fmt:
	gofumpt -w .

fmt-check:
	@output="$$(gofumpt -l .)"; \
	if [ -n "$$output" ]; then \
		echo "Files not formatted with gofumpt:"; \
		echo "$$output"; \
		exit 1; \
	fi

ci: fmt-check vet lint test
