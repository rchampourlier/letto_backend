package js

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/spf13/afero"

	"gitlab.com/letto/letto_backend/exec/values"
	"gitlab.com/letto/letto_backend/services/events"
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

// Execute runs the JS execution environment for the specified
// event.
//
// The event file is created in `exec/js/src`. This file
// provides the container with the event (e.g. request's headers,
// body...).
//
// The container image is built from the `exec/js` directory which
// contains the `Dockerfile` and the `src` directory.
func (r *Runner) Execute(event events.ActivatedTriggerEvent) (values.Output, error) {
	output := values.Output{}

	// Dump the context to a file passed to the script
	eventJS, err := generateEventJS(event.Data())
	if err != nil {
		return output, err
	}

	eventFileName := "event-" + util.Timestamp(time.Now()) + ".js"
	eventFilePath := path.Join(r.appDataDir, eventFileName)
	err = afero.WriteFile(r.Fs, eventFilePath, []byte(eventJS), 0777)
	if err != nil {
		return output, err
	}

	output, err = dockerExecute(eventFileName, r.hostDataDir, r.execDataDir)
	if err != nil {
		fmt.Printf("Error while running JS workflows: %s\n", err)
		return output, err
	}

	err = r.Fs.Remove(eventFilePath)
	if err != nil {
		fmt.Printf("Could not remove context temp file: %s\n", eventFilePath)
	}

	return output, nil
}

// generateEventJS generates JS script to provide the event's data to
// the JS scripts.
func generateEventJS(eventData events.EventData) (string, error) {
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return "", err
	}

	// Insert context data in the template JS script
	eventJS := strings.Replace(eventJSTemplate, "{{event}}", string(eventJSON), 1)
	return eventJS, nil
}

var eventJSTemplate = `
var event = {{event}};
module.exports = event;
`
