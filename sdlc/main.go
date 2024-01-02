package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
)

var (
	rootCmd = cobra.Command{
		Use:                "sdlc",
		PersistentPreRunE:  connect,
		PersistentPostRunE: disconnect,
		SilenceErrors:      true,
	}

	client *dagger.Client
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	assert(rootCmd.ExecuteContext(ctx))
}

func assert(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func connect(cmd *cobra.Command, args []string) error {
	var err error

	client, err = dagger.Connect(cmd.Context(), dagger.WithLogOutput(os.Stderr))
	// dagger.WithConn(WithRemoteEngineConn("aaa"))

	return err
}

func disconnect(cmd *cobra.Command, args []string) error {
	return client.Close()
}
