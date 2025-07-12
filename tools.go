//go:build tools
// +build tools

package anilistarr

import (
	_ "github.com/Khan/genqlient"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "github.com/wadey/gocovmerge"
	_ "golang.org/x/tools/cmd/stringer"
)
