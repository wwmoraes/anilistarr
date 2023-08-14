package anilistarr

//go:generate go run ./cmd/version/... -package telemetry -name anilistarr -namespace api -version "$VERSION" -output internal/telemetry/constants.go

//// DISABLED: this generator needs a db to derive the code from
// go:generate xo schema "file:tmp/media.db?loc=auto" -o internal/drivers/stores/models

//// TODO: genqlient generator
//// cd internal/drivers/anilist && genqlient
