// Package stores provides [usecases.Store] implementations for various
// different storage technologies.
package main

const (
	// StoreBadger represents BadgerDB-backed store drivers
	StoreBadger StoreType = "badger"
	// StoreSQL represents SQL/RDBMS-backed store drivers
	StoreSQL StoreType = "sql"
)

// StoreType represents a type of store driver backing technology
type StoreType string
