package stabber

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/pkg/jsonmessage"
)

type ImageLoader interface {
	ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error)
}

func (stabber *Stabber) LoadDockerIgnore(fsys fs.FS, name string) ([]string, error) {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, err
	}

	// normalize EOLs
	data = bytes.ReplaceAll(data, []byte{'\r'}, []byte{'\n'})
	data = bytes.ReplaceAll(data, []byte{'\n', '\n'}, []byte{'\n'})

	entries := make([]string, 0, bytes.Count(data, []byte{'\n'}))
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		line = bytes.Trim(line, " ")

		if bytes.Equal(line, []byte{}) {
			continue
		}

		if bytes.HasPrefix(line, []byte{'#'}) {
			continue
		}

		entries = append(entries, string(line))
	}

	return entries, nil
}

func (stabber *Stabber) LoadDockerImage(ctx context.Context, container *dagger.Container, loader ImageLoader) (string, error) {
	tmpDir, err := os.MkdirTemp("", "stabber-load-image.*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "container.tar")

	_, err = container.AsTarball().Export(ctx, tarPath)
	if err != nil {
		return "", err
	}

	fd, err := os.OpenFile(tarPath, os.O_RDONLY, 0640)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	res, err := loader.ImageLoad(ctx, fd, true)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	msg := jsonmessage.JSONMessage{}
	json.Unmarshal(data, &msg)

	// absolutely beautiful API Docker, thank you...
	imageID := strings.TrimPrefix(strings.Trim(msg.Stream, "\r\n"), "Loaded image ID: ")

	return imageID, nil
}

func (stabber *Stabber) DockerBuildContext(dir string) (*dagger.Directory, error) {
	excludes, err := stabber.LoadDockerIgnore(os.DirFS(dir), ".dockerignore")
	if errors.Is(err, fs.ErrNotExist) {
		excludes = []string{}
	} else if err != nil {
		return nil, err
	}

	return stabber.Client.Host().
		Directory(dir, dagger.HostDirectoryOpts{
			Exclude: excludes,
		}).
		WithFile("Dockerfile", stabber.Client.Host().File("Dockerfile")), nil
}

func (stabber *Stabber) DockerLintContainer(ctx context.Context) (*dagger.Container, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dockerfilePath := filepath.Join(pwd, "Dockerfile")
	dockerfile := stabber.Client.Host().File("Dockerfile")

	hadolint := stabber.Client.Container().From(ImageHadolint).
		WithMountedFile(dockerfilePath, dockerfile).
		WithExec([]string{"hadolint", dockerfilePath})

	// TODO check also $XDG_CONFIG_HOME/hadolint.yaml and ~/.config/hadolint.yaml
	_, err = os.Stat(".hadolint.yaml")
	if err == nil {
		fmt.Println("adding hadoling config")
		hadolintConfig := stabber.Client.Host().File(".hadolint.yaml")
		hadolint = hadolint.WithMountedFile("/root/.config/hadolint.yaml", hadolintConfig)
	}

	return hadolint.Sync(ctx)
}
