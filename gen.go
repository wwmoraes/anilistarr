package anilistarr

// SQL models and queries
//go:generate sqlc generate

// API server
//go:generate oapi-codegen -generate types,chi-server,spec -package api -o internal/api/api.gen.go swagger.yaml
// Thank you oapi-codegen...
//go:generate sh -c "sed -i'' '/var err error/d' internal/api/api.gen.go"
