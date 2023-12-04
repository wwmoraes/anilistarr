package stores

import (
	"github.com/wwmoraes/anilistarr/internal/drivers/persistence"
)

type BadgerOptions = persistence.BadgerOptions
type BadgerLogr = persistence.BadgerLogr

var (
	NewBadger = persistence.NewBadger
)
