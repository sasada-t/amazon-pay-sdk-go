.PHONY: http.recurring
http.recurring:
	@go run github.com/air-verse/air -c .air_recurring.toml

.PHONY: build
build: build.recurring

.PHONY: build.recurring
build.recurring:
	@go build -o ./.bin/recurring ./example/recurring/main.go

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