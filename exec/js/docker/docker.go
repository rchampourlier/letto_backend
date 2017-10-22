// This file provides JS execution through Docker.
//
// It comes with an unique `Run` function that may
// be used to simply execute a Docker container
// with a reduced set of parameters.

package docker

import (
	"archive/tar"
	"bytes"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"

	"github.com/rchampourlier/letto_go/exec/values"
)

const imageNameAndTag = "letto/exec-js:latest"

// Prepare prepares the Docker environment. In particular, it builds
// the container image.
//
// `rootDir` specified the root directory that will be used to build
// the Docker image.
//
// `out` is used to write logs generated by Docker while building
// the image.
func Prepare(rootDir string, out io.Writer) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	options := types.ImageBuildOptions{
		Tags: []string{imageNameAndTag},
	}
	var resp types.ImageBuildResponse

	// buildContext is a io.Reader with a tar archive containing the
	// Dockerfile.
	buildContext, err := makeBuildContext(rootDir)
	if err != nil {
		return err
	}

	resp, err = cli.ImageBuild(ctx, buildContext, options)
	if err != nil {
		return err
	}

	io.Copy(out, resp.Body)
	return nil
}

// Run executes a Docker container with the specified
// parameters.
func Run(srcDir string, contextFileName string) (values.Output, error) {

	output := values.Output{}
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return output, err
	}

	// The image is created locally so we don't need to pull
	// it anymore.
	/*_, err = cli.ImagePull(ctx, imageNameAndTag, types.ImagePullOptions{})
	if err != nil {
		return output, err
	}*/

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      imageNameAndTag,
		Cmd:        []string{"node", "./main.js", "./" + contextFileName},
		Volumes:    map[string]struct{}{"/usr/src/app": {}},
		WorkingDir: "/usr/src/app",
	}, &container.HostConfig{
		Binds: []string{srcDir + ":/usr/src/app"},
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

func makeBuildContext(dirPath string) (io.Reader, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new tar archive.
	tw := tar.NewWriter(buf)

	// Add the directory's contents to the archive
	tarDirectory(dirPath, tw)

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func tarFile(filePath string, tw *tar.Writer, fi os.FileInfo) error {
	fr, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fr.Close()

	hdr := &tar.Header{
		Name: fi.Name(),
		Mode: 0600,
		Size: fi.Size(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := io.Copy(tw, fr); err != nil {
		return err
	}

	return nil
}

func tarDirectory(dirPath string, tw *tar.Writer) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	fileinfos, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range fileinfos {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			tarDirectory(curPath, tw)
		} else {
			tarFile(curPath, tw, fi)
		}
	}

	return nil
}