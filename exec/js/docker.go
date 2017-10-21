// This file provides JS execution through Docker.
//
// It comes with an unique `Run` function that may
// be used to simply execute a Docker container
// with a reduced set of parameters.

package js

import (
	"bytes"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// DockerConfig defines the config of the Docker container
// to be run.
type DockerConfig struct {
	Image      string
	Command    []string
	Volumes    map[string]struct{}
	WorkingDir string
	Binds      []string
}

// Run executes a Docker container with the specified
// parameters.
func Run(cfg DockerConfig) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ImagePull(ctx, cfg.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      cfg.Image,
		Cmd:        cfg.Command,
		Volumes:    cfg.Volumes,
		WorkingDir: cfg.WorkingDir,
	}, &container.HostConfig{
		Binds: cfg.Binds,
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	//err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	io.Copy(os.Stdout, out)
}
