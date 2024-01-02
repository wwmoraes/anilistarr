package main

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	moby "github.com/moby/moby/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/wwmoraes/anilistarr/sdlc/stabber"
)

var buildFlags *pflag.FlagSet

func init() {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build binary of application",
		RunE:  dagger2CobraCmd(build),
	}

	buildFlags = cmd.Flags()

	cmd.Flags().Bool("load", false, "load the built container on the host engine")
	cmd.Flags().String("version", "0.0.0", "load the built container on the host engine")

	rootCmd.AddCommand(cmd)
}

func build(ctx context.Context, client *dagger.Client, options Options) error {
	client = client.Pipeline("build", dagger.PipelineOpts{
		Description: "builds the application binary",
	})

	stab := &stabber.Stabber{Client: client}

	modCache := client.CacheVolume("go-mod")
	buildCache := client.CacheVolume("go-build")

	rootDir := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Include: []string{
			"cmd/",
			"internal/",
			".dockerignore",
			"Dockerfile",
			"go.mod",
			"go.sum",
			"*.go",
		},
	})

	buildContext, err := stab.DockerBuildContext(".")
	if err != nil {
		return err
	}

	version, err := options.GetString("version")
	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	tags, err := repo.Tags()
	if err != nil {
		return err
	}

	tag, err := tags.Next()
	if err != nil {
		return err
	}

	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	revlist.Objects(repo.Storer, []plumbing.Hash{
		tag.Hash(),
		headRef.Hash(),
	}, []plumbing.Hash{})

	tagName := tag.Name().Short()
	if tag.Hash() != headRef.Hash() {
		tagName = fmt.Sprintf("%s-nightly.%s", tagName, headRef.Hash().String())
	}

	container := rootDir.DockerBuild().
		WithMountedCache("/go/pkg/mod", modCache).
		WithMountedCache("/root/.cache/go-build", buildCache).
		Build(buildContext, dagger.ContainerBuildOpts{
			BuildArgs: []dagger.BuildArg{
				{
					Name:  "VERSION",
					Value: version,
				},
			},
		})

	container.Sync(ctx)

	load, err := options.GetBool("load")
	if err != nil {
		return err
	}

	if !load {
		return nil
	}

	// load the built image on the host
	dockerClient, err := moby.NewClientWithOpts(moby.FromEnv, moby.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	imageID, err := stab.LoadDockerImage(ctx, container, dockerClient)
	if err != nil {
		return err
	}

	err = dockerClient.ImageTag(ctx, imageID, "wwmoraes/anilistarr:latest")
	if err != nil {
		return err
	}

	err = dockerClient.ImageTag(ctx, imageID, fmt.Sprintf("wwmoraes/anilistarr:%s", tagName))
	if err != nil {
		return err
	}

	return nil
}
