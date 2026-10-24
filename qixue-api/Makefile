.PHONY: fmt vet test build check

fmt:
	gofmt -w .
	@command -v goimports >/dev/null 2>&1 && goimports -w -local qixue-api . || true

vet:
	go vet ./...

test:
	go test ./...

build:
	go build -o bin/server ./cmd/server

check: vet test
