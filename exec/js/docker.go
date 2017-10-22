// This file provides JS execution through Docker.
//
// It comes with an unique `Run` function that may
// be used to simply execute a Docker container
// with a reduced set of parameters.

package js

import (
	"bytes"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"

	"github.com/rchampourlier/letto_go/exec/values"
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
func Run(cfg DockerConfig) (values.Output, error) {
	output := values.Output{}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return output, err
	}

	_, err = cli.ImagePull(ctx, cfg.Image, types.ImagePullOptions{})
	if err != nil {
		return output, err
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
		return output, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return output, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return output, err
		}
	case <-statusCh:
	}

	logs, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     false,
		Details:    false,
	})
	if err != nil {
		return output, err
	}

	// As explained in github.com/moby/moby/client/container_logs.go
	// we use `stdcopy.Stdcopy` to demulitplex the logs.
	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)
	stdcopy.StdCopy(bufOut, bufErr, logs)
	output.Stdout = bufOut.String()
	output.Stderr = bufErr.String()

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return output, err
	}

	return output, nil
}
