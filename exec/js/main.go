package main

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ImagePull(ctx, "node", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	dir, _ := os.Getwd()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      "node:4",
		Cmd:        []string{"node", "./test.js"},
		Volumes:    map[string]struct{}{"/usr/src/app": {}},
		WorkingDir: "/usr/src/app",
	}, &container.HostConfig{
		Binds: []string{dir + ":/usr/src/app"},
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

	io.Copy(os.Stdout, out)
}
