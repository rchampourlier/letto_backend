package js

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/spf13/afero"

	"gitlab.com/letto/letto_backend/exec/values"
	"gitlab.com/letto/letto_backend/util"
)

var tmpDirPrefix = "letto"

// Runner is the structure holding configuration for
// the javascript execution environment.
//
// `hostDataDir` is where the data is contained on the host
// `execDataDir` is where the data will be contained on the
//   container.
type Runner struct {
	Fs          afero.Fs
	hostDataDir string
	appDataDir  string
	execDataDir string
}

// NewRunner creates a new JsRunner
func NewRunner(fs afero.Fs, hostDataDir string, appDataDir string, execDataDir string) (Runner, error) {
	r := Runner{
		Fs:          fs,
		hostDataDir: hostDataDir,
		appDataDir:  appDataDir,
		execDataDir: execDataDir,
	}
	// Disabled because the Docker image is now prepared
	// using docker-compose.
	/*if err := dockerPrepare(os.Stdout); err != nil {
		return r, err
	}*/
	return r, nil
}

// Execute runs the JS execution environment with the specified
// group and context.
//
// The `group` determines which workflows are run (all workflows
// defined by scripts under `src/workflows/<group>`) and each
// workflow is passed the context.
//
// The context file is created in `exec/js/src`. This file
// provides the container with the context (e.g. request's headers,
// body...).
//
// The container image is built from the `exec/js` directory which
// contains the `Dockerfile` and the `src` directory.
func (r *Runner) Execute(group string, ctx values.Context) (values.Output, error) {
	output := values.Output{}

	// Dump the context to a file passed to the script
	contextJS, err := generateContextJS(ctx)
	if err != nil {
		return output, err
	}

	contextFileName := "context-" + util.Timestamp(time.Now()) + ".js"
	contextFilePath := path.Join(r.appDataDir, contextFileName)
	err = afero.WriteFile(r.Fs, contextFilePath, []byte(contextJS), 0777)
	if err != nil {
		return output, err
	}

	output, err = dockerExecute(contextFileName, r.hostDataDir, r.execDataDir)
	if err != nil {
		fmt.Printf("Error while running JS workflows: %s\n", err)
		return output, err
	}

	err = r.Fs.Remove(contextFilePath)
	if err != nil {
		fmt.Printf("Could not remove context temp file: %s\n", contextFilePath)
	}

	return output, nil
}

// contextJS generates JS script to provide the context to
// the JS scripts.
func generateContextJS(ctx values.Context) (string, error) {
	contextDataJS, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}

	// Insert context data in the template JS script
	contextJS := strings.Replace(contextJSTemplate, "{{context}}", string(contextDataJS), 1)
	return contextJS, nil
}

var contextJSTemplate = `
var context = {{context}};
module.exports = context;
`
