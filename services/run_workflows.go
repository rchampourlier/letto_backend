package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/events"
	"github.com/rchampourlier/letto_go/exec"
)

var tmpDirPrefix = "letto"

// RunWorkflows stores the service's config.
type RunWorkflows struct {
	Fs afero.Fs
}

// NewRunWorkflows creates a new RunWorkflowsService with
// the specified `afero.Fs` filesystem abstraction.
func NewRunWorkflows(fs afero.Fs) *RunWorkflows {
	return &RunWorkflows{fs}
}

// OnReceivedWebhook runs all workflows for the specified
// `group` with the specified event.
//
// TODO: refactor by extracting JS-related code to exec/js
func (s *RunWorkflows) OnReceivedWebhook(event events.ReceivedWebhook) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create a tmp dir mounted in the container to provide the scripts
	// and data.
	tmpDir, err := ioutil.TempDir("", tmpDirPrefix)
	if err != nil {
		return err
	}

	// Copy the `main.js` script that will load the workflows, request
	// and data, and run the workflows.

	// Copy the workflow scripts, selecting only the ones from
	// the `group`.
	//group := ctx.Group
	// TODO: selection by group
	s.copyFile(path.Join(rootDir, "exec", "js", "main.js"), path.Join(tmpDir, "main.js"))

	// Dump the request's data to a file passed to the script
	requestJS, err := generateRequestJS(event.Body, event.Headers)
	if err != nil {
		return err
	}
	err = afero.WriteFile(s.Fs, path.Join(tmpDir, "request.js"), []byte(requestJS), 0777)
	if err != nil {
		return err
	}
	s.copyFile(path.Join(rootDir, "exec", "js", "data.js"), path.Join(tmpDir, "data.js"))
	s.copyFile(path.Join(rootDir, "credentials.js"), path.Join(tmpDir, "credentials.js"))

	exec.RunJS(tmpDir, "./main.js")

	err = s.Fs.RemoveAll(tmpDir)
	if err != nil {
		fmt.Printf("Could not remove tmp dir: %s\n", tmpDir)
	}

	return nil
}

// TODO: handle requests where body is not JSON
func generateRequestJS(body string, headers map[string][]string) (string, error) {
	var requestJS string

	if len(body) > 0 {
		// TODO: test that body is JSON
		var bodyParsed map[string]interface{}
		err := json.Unmarshal([]byte(body), &bodyParsed)
		if err != nil {
			return "", err
		}

		bodyAsJSON, err := json.Marshal(bodyParsed)
		if err != nil {
			return "", err
		}
		requestJS = strings.Replace(requestJSTemplate, "{{body}}", string(bodyAsJSON), 1)
	} else {
		requestJS = strings.Replace(requestJSTemplate, "{{body}}", "null", 1)
	}

	headersAsJSON, err := json.Marshal(headers)
	if err != nil {
		return "", err
	}

	// Insert request data in the template JS script
	requestJS = strings.Replace(requestJS, "{{headers}}", string(headersAsJSON), 1)
	return requestJS, nil
}

var requestJSTemplate = `
var body = {{body}};
var headers = {{headers}};
module.exports = {
	body: body, 
	headers: headers
};
`

func (s *RunWorkflows) copyFile(src string, dst string) error {
	data, err := afero.ReadFile(s.Fs, src)
	if err != nil {
		return err
	}
	err = afero.WriteFile(s.Fs, dst, data, 0777)
	if err != nil {
		return err
	}
	return nil
}
