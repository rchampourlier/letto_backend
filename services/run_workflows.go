package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/exec"
)

var appFs = afero.NewOsFs()
var tmpDirPrefix = "letto"

// WorkflowData is a structure that is used to
// pass data to execute the workflow with.
type WorkflowData struct {
	RequestBody    string
	RequestHeaders map[string][]string
}

// RunWorkflows run all workflows for the specified `group`
// with the specified data.
//
// TODO: refactor by extracting JS-related code to exec/js
func RunWorkflows(group string, data WorkflowData) error {
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
	copyFile(path.Join(rootDir, "exec", "js", "main.js"), path.Join(tmpDir, "main.js"))

	// Dump the request's data to a file passed to the script
	requestJS, err := generateRequestJS(data)
	if err != nil {
		return err
	}
	err = afero.WriteFile(appFs, path.Join(tmpDir, "request.js"), []byte(requestJS), 0777)
	if err != nil {
		return err
	}
	copyFile(path.Join(rootDir, "exec", "js", "data.js"), path.Join(tmpDir, "data.js"))
	copyFile(path.Join(rootDir, "credentials.js"), path.Join(tmpDir, "credentials.js"))

	exec.RunJS(tmpDir, "./main.js")

	err = appFs.RemoveAll(tmpDir)
	if err != nil {
		fmt.Printf("Could not remove tmp dir: %s\n", tmpDir)
	}

	return nil
}

// TODO: handle requests where body is not JSON
func generateRequestJS(data WorkflowData) (string, error) {
	body := data.RequestBody
	headers := data.RequestHeaders
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

func copyFile(src string, dst string) error {
	data, err := afero.ReadFile(appFs, src)
	if err != nil {
		return err
	}
	err = afero.WriteFile(appFs, dst, data, 0777)
	if err != nil {
		return err
	}
	return nil
}
