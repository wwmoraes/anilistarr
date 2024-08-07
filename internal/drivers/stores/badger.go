package stores

import (
	"github.com/wwmoraes/anilistarr/internal/drivers/persistence"
)

type (
	BadgerOptions = persistence.BadgerOptions
	BadgerLogr    = persistence.BadgerLogr
)

var NewBadger = persistence.NewBadger
