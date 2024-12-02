package main

const (
	// CacheBadger represents BadgerDB-backed cache drivers
	CacheBadger CacheType = "badger"
	// CacheBolt represents BoltDB-backed cache drivers
	CacheBolt CacheType = "bolt"
	// CacheRedis represents Redis-backed cache drivers
	CacheRedis CacheType = "redis"
	// CacheSQL represents SQL/RDBMS-backed cache drivers
	CacheSQL CacheType = "sql"
)

// CacheType represents a type of cache-backing technology
type CacheType string
