.PHONY: build
build:
	@go build -o ./.bin/main ./example/recurring/main.go

.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix
	$(call format)

.PHONY: format
format:
	$(call format)

define format
	@go fmt ./... 
	@go run golang.org/x/tools/cmd/goimports -w ./ 
	@go run mvdan.cc/gofumpt -l -w .
	@go mod tidy
endef