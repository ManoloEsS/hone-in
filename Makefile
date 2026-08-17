GO ?= go
GOFMT ?= gofmt
DB_PATH ?= ./hone-in.db
HTTP_ADDR ?= :8080

.PHONY: check fmt-check run test vet

check: fmt-check vet test

run:
	DB_PATH="$(DB_PATH)" HTTP_ADDR="$(HTTP_ADDR)" $(GO) run ./cmd/server

fmt-check:
	@files="$$(find . -type f -name '*.go' -not -path './.jj/*')"; \
	if [ -n "$$files" ]; then \
		test -z "$$($(GOFMT) -l $$files)"; \
	fi

test:
	@files="$$(find . -type f -name '*.go' -not -path './.jj/*')"; \
	if [ -n "$$files" ]; then \
		$(GO) test ./...; \
	else \
		echo "no Go source files; skipping tests"; \
	fi

vet:
	@files="$$(find . -type f -name '*.go' -not -path './.jj/*')"; \
	if [ -n "$$files" ]; then \
		$(GO) vet ./...; \
	else \
		echo "no Go source files; skipping vet"; \
	fi
