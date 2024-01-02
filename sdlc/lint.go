package main

import (
	"context"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/wwmoraes/anilistarr/sdlc/stabber"
)

var lintFlags *pflag.FlagSet

func init() {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "execute linters on the source code",
		RunE:  dagger2CobraCmd(lint),
	}

	lintFlags = cmd.Flags()

	rootCmd.AddCommand(cmd)
}

func lint(ctx context.Context, client *dagger.Client, options Options) error {
	client = client.Pipeline("lint", dagger.PipelineOpts{
		Description: "execute linters on the source code",
	})

	stab := stabber.Stabber{Client: client}

	_, err := stab.DockerLintContainer(ctx)

	return stabber.TryUnwrapExecError(err)
}
