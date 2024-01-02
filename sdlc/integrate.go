package main

import (
	"context"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	cmd := &cobra.Command{
		Use:   "integrate",
		Short: "validates changes to the repository",
		RunE:  dagger2CobraCmd(integrate),
	}

	cmd.Flags().AddFlagSet(lintFlags)
	cmd.Flags().AddFlagSet(buildFlags)
	cmd.Flags().AddFlagSet(testFlags)

	rootCmd.AddCommand(cmd)
}

func integrate(ctx context.Context, client *dagger.Client, options Options) error {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return lint(gCtx, client, options)
	})

	group.Go(func() error {
		return build(gCtx, client, options)
	})

	group.Go(func() error {
		return test(gCtx, client, options)
	})

	return group.Wait()
}
