package main

import (
	"context"
	"net/http"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

type Options interface {
	// GetBool return the bool value of a flag with the given name
	GetBool(name string) (bool, error)
	// GetString return the string value of a flag with the given name
	GetString(name string) (string, error)
}

type DaggerRunE func(ctx context.Context, client *dagger.Client, options Options) error
type CobraRunE func(cmd *cobra.Command, args []string) error

func dagger2CobraCmd(block DaggerRunE) CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		return block(cmd.Context(), client, cmd.Flags())
	}
}

func WithRemoteEngineConn(host string) *RemoteEngineConn {
	return &RemoteEngineConn{
		Client: *http.DefaultClient,
		host:   host,
	}
}

type RemoteEngineConn struct {
	http.Client

	host string
}

func (conn *RemoteEngineConn) Host() string {
	return conn.host
}

func (conn *RemoteEngineConn) Close() error {
	return nil
}
