package anilistarr

// SQL models and queries
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

// API server
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,chi-server,spec -package api -o internal/api/api.gen.go swagger.yaml
// Thank you oapi-codegen...
//go:generate sed --in-place= "/var err error/d" internal/api/api.gen.go
