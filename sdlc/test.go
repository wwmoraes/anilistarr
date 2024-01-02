package main

import (
	"context"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var testFlags *pflag.FlagSet

func init() {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "executes tests of the application",
		RunE:  dagger2CobraCmd(test),
	}

	testFlags = cmd.Flags()

	rootCmd.AddCommand(cmd)
}

func test(ctx context.Context, client *dagger.Client, options Options) error {
	client = client.Pipeline("test")

	modCache := client.CacheVolume("go-mod")
	buildCache := client.CacheVolume("go-build")
	rootDir := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Include: []string{
			"cmd/",
			"internal/",
			"*.go",
			".dockerignore",
			"Dockerfile",
			"go.mod",
			"go.sum",
		},
	})

	container := rootDir.DockerBuild()
	container.WithMountedCache("/go/pkg/mod", modCache)
	container.WithMountedCache("/root/.cache/go-build", buildCache)

	_, err := container.Build(rootDir, dagger.ContainerBuildOpts{
		BuildArgs: []dagger.BuildArg{
			{
				Name:  "VERSION",
				Value: "0.0.0",
			},
		},
		Target: "build",
	}).WithExec([]string{"go", "test", "-race -v", "./..."}, dagger.ContainerWithExecOpts{
		SkipEntrypoint: true,
	}).Stdout(ctx)
	if err != nil {
		return err
	}

	return nil
}
