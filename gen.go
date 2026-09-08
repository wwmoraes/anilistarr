package anilistarr

// SQL models and queries
//go:generate go tool sqlc generate

// API server
//go:generate go tool oapi-codegen -generate types,chi-server,spec -package api -o internal/api/api.gen.go swagger.yaml
// Thank you oapi-codegen...
//go:generate sed --in-place= "/var err error/d" internal/api/api.gen.go
//go:generate sed --in-place= "/_ = err/d" internal/api/api.gen.go
