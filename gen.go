package anilistarr

// DISABLED: this generator needs a db to derive the code from
// echo go:generate xo schema "file:tmp/media.db?loc=auto" -o internal/drivers/stores/models

// API server
//go:generate oapi-codegen -generate types,chi-server,spec -package api -o internal/api/api.gen.go swagger.yaml

// Thank you oapi-codegen...
//go:generate sh -c "sed -i'' '/var err error/d' internal/api/api.gen.go"
