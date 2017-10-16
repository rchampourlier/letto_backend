package exec

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/exec/js"
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

// Execute executes the specified script in a JS execution
// environment. The script is passed the `body` string.
//
// `dir` is expected to be the path to a local directory
// which contains the JS code to be executed. The container
// will run the `main.js` file present in this directory.
//
// `main` is the path, local to `dir`, of the main JS file
// to be run by the container.
func (runner *JsRunner) Execute(credentialsPath string, group string, ctx Context) error {

	// Copy the `main.js` script that will load the workflows, request
	// and data, and run the workflows.
	rootDir, err := os.Getwd()

	// Create a tmp dir mounted in the container to provide the scripts
	// and data.
	tmpDir, err := ioutil.TempDir("", tmpDirPrefix)
	if err != nil {
		return err
	}

	runner.copyFile(credentialsPath, path.Join(tmpDir, "credentials.js"))

	// Copy the workflow scripts, selecting only the ones from
	// the `group`.
	// TODO: selection by group
	runner.copyFile(path.Join(rootDir, "exec", "js", "main.js"), path.Join(tmpDir, "main.js"))

	// Dump the context to a file passed to the script
	contextJS, err := generateContextJS(ctx)
	if err != nil {
		return err
	}
	err = afero.WriteFile(runner.Fs, path.Join(tmpDir, "context.js"), []byte(contextJS), 0777)
	if err != nil {
		return err
	}
	runner.copyFile(path.Join(rootDir, "exec", "js", "data.js"), path.Join(tmpDir, "data.js"))

	cfg := config(tmpDir, "./main.js")
	js.Run(cfg)

	err = runner.Fs.RemoveAll(tmpDir)
	if err != nil {
		fmt.Printf("Could not remove tmp dir: %s\n", tmpDir)
	}

	return nil
}

func config(dir string, main string) js.DockerConfig {
	var cfg = js.DockerConfig{
		Image:      "node:latest",
		Command:    []string{"node", main},
		Volumes:    map[string]struct{}{"/usr/src/app": {}},
		WorkingDir: "/usr/src/app",
		Binds:      []string{dir + ":/usr/src/app"},
	}
	return cfg
}

// contextJS generates JS script to provide the context to
// the JS scripts.
func generateContextJS(ctx Context) (string, error) {
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

func (runner *JsRunner) copyFile(src string, dst string) error {
	data, err := afero.ReadFile(runner.Fs, src)
	if err != nil {
		return err
	}
	err = afero.WriteFile(runner.Fs, dst, data, 0777)
	if err != nil {
		return err
	}
	return nil
}
