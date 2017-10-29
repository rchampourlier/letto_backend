package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/afero"

	"gitlab.com/letto/letto_backend/exec/js/docker"
	"gitlab.com/letto/letto_backend/exec/values"
	"gitlab.com/letto/letto_backend/util"
)

var tmpDirPrefix = "letto"

// JsRunner is the structure holding configuration for
// the javascript execution environment.
type JsRunner struct {
	Fs afero.Fs
}

// NewJsRunner creates a new JsRunner
func NewJsRunner(fs afero.Fs) JsRunner {
	return JsRunner{
		Fs: fs,
	}
}

// Prepare prepares the JS execution environment with the
// specified root dir.
func (runner *JsRunner) Prepare(rootDir string) error {
	err := docker.Prepare(rootDir, os.Stdout)
	return err
}

// Execute runs the JS execution environment with the specified
// group and context.
//
// The `group` determines which workflows are run (all workflows
// defined by scripts under `src/workflows/<group>`) and each
// workflow is passed the context.
//
// A temporary context file is created in `exec/js/src`. This file
// provides the container with the context (e.g. request's headers,
// body...).
//
// The container then mounts the `exec/js/src` directory which
// also contains the JS source code (`main.js` and the workflows)
// that will be executed when running the container.
func (runner *JsRunner) Execute(group string, ctx values.Context) (values.Output, error) {
	rootDir, err := os.Getwd()
	output := values.Output{}

	// `jsSrcDir` is mounted on the container
	jsSrcDir := path.Join(rootDir, "exec", "js", "src")

	// Dump the context to a file passed to the script
	contextJS, err := generateContextJS(ctx)
	if err != nil {
		return output, err
	}

	contextFileName := "context-" + util.Timestamp(time.Now()) + ".js"
	contextJSPath := path.Join(jsSrcDir, contextFileName)
	err = afero.WriteFile(runner.Fs, contextJSPath, []byte(contextJS), 0777)
	if err != nil {
		return output, err
	}

	output, err = docker.Run(jsSrcDir, contextFileName)
	if err != nil {
		fmt.Printf("Error while running JS workflows: %s\n", err)
	}

	err = runner.Fs.Remove(contextJSPath)
	if err != nil {
		fmt.Printf("Could not remove context temp file: %s\n", contextJSPath)
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
