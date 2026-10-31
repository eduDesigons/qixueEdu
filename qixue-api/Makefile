.PHONY: fmt vet test build check tidy

fmt:
	gofmt -w .
	@command -v goimports >/dev/null 2>&1 && goimports -w -local qixue-api . || true

vet:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

build:
	go build -o bin/server ./cmd/server

check: vet test
