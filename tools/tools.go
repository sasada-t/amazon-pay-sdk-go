//go:build tools

package tool

import (
	_ "github.com/air-verse/air"                            //nolint
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" //nolint
	_ "golang.org/x/tools/cmd/goimports"                    //nolint
	_ "mvdan.cc/gofumpt"                                    //nolint
)
