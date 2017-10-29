package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/afero"

	"gitlab.com/letto/letto_backend/exec/js"
	"gitlab.com/letto/letto_backend/exec/values"
	"gitlab.com/letto/letto_backend/util"
)

var tmpDirPrefix = "letto"

// PrepareJsRunner prepares the JS execution environment with the
// specified root dir. It must be called before running
// `Execute` on a `JsRunner`.
func PrepareJsRunner(rootDir string) error {
	err := js.Prepare(rootDir, os.Stdout)
	return err
}

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
func (runner *JsRunner) Execute(group string, ctx values.Context) (values.Output, error) {
	output := values.Output{}

	// TODO: should be injected with the location of the
	// `exec/js` directory.
	cwd, err := os.Getwd()
	if err != nil {
		return output, err
	}

	// Dump the context to a file passed to the script
	contextJS, err := generateContextJS(ctx)
	if err != nil {
		return output, err
	}

	contextFileName := "context-" + util.Timestamp(time.Now()) + ".js"
	contextFilePath := path.Join(cwd, "exec", "js", "src", contextFileName)
	// TODO: JsRunner should not have to know to put the
	//   JsContextFile inside `exec/js`.
	err = afero.WriteFile(runner.Fs, contextFilePath, []byte(contextJS), 0777)
	if err != nil {
		return output, err
	}

	output, err = js.Run(contextFileName)
	if err != nil {
		fmt.Printf("Error while running JS workflows: %s\n", err)
	}

	err = runner.Fs.Remove(contextFilePath)
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
