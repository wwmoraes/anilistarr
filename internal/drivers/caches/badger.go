package caches

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-logr/logr"
	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
)

type badgerCache struct {
	*badger.DB
}

func NewBadger(path string, options ...BadgerOption) (adapters.Cache, error) {
	opt := badger.DefaultOptions(path)
	for _, optFn := range options {
		opt = optFn.Apply(opt)
	}

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return &badgerCache{db}, nil
}

func (c *badgerCache) GetString(ctx context.Context, key string) (string, error) {
	_, span := telemetry.StartFunction(ctx)
	defer span.End()

	var value string

	err := c.View(func(txn *badger.Txn) error {
		data, err := txn.Get([]byte(key))
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		if data.IsDeletedOrExpired() {
			return nil
		}

		return data.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
	})

	return value, span.Assert(err)
}

func (c *badgerCache) SetString(ctx context.Context, key, value string, options ...adapters.CacheOption) error {
	_, span := telemetry.StartFunction(ctx)
	defer span.End()

	params, err := adapters.NewCacheParams(options...)
	if err != nil {
		return err
	}

	return span.Assert(c.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), []byte(value))

		return txn.SetEntry(entry.WithTTL(params.TTL))
	}))
}

type BadgerOption interface {
	Apply(badger.Options) badger.Options
}

type BadgerOptionFn func(badger.Options) badger.Options

func (fn BadgerOptionFn) Apply(options badger.Options) badger.Options {
	return fn(options)
}

func WithLogger(log logr.Logger) BadgerOption {
	return BadgerOptionFn(func(options badger.Options) badger.Options {
		return options.WithLogger(&badgerLogger{log})
	})
}

func WithInMemory(b bool) BadgerOption {
	return BadgerOptionFn(func(option badger.Options) badger.Options {
		return option.WithInMemory(b)
	})
}

type badgerLogger struct {
	log logr.Logger
}

func (log *badgerLogger) Errorf(format string, a ...interface{}) {
	log.log.Error(fmt.Errorf(format, a...), "Badger Error")
}

func (log *badgerLogger) Warningf(format string, a ...interface{}) {
	log.log.Info(fmt.Sprintf(format, a...), "Badger Warning")
}

func (log *badgerLogger) Infof(format string, a ...interface{}) {
	log.log.Info(fmt.Sprintf(format, a...), "Badger Info")
}

func (log *badgerLogger) Debugf(format string, a ...interface{}) {
	log.log.Info(fmt.Sprintf(format, a...), "Badger Debug")
}
