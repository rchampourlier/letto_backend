package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/goadesign/goa"
	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/app"
	"github.com/rchampourlier/letto_go/exec"
)

var tmpDirPrefix = "letto"
var appFs = afero.NewOsFs()

// TriggersController implements the triggers resource.
type TriggersController struct {
	*goa.Controller
}

// NewTriggersController creates a triggers controller.
func NewTriggersController(service *goa.Service) *TriggersController {
	return &TriggersController{Controller: service.NewController("TriggersController")}
}

// Webhook runs the webhook action.
//
// TODO: refactor by extracting JS execution logic
//   to exec/js.go
func (c *TriggersController) Webhook(ctx *app.WebhookTriggersContext) error {
	// TriggersController_Webhook: start_implement

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
	requestJS, err := generateRequestJS(ctx)
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

	// TriggersController_Webhook: end_implement
	return nil
}

// TODO: handle requests where body is not JSON
func generateRequestJS(ctx *app.WebhookTriggersContext) (string, error) {
	var body string
	var headers map[string][]string
	var requestJS string

	body = readBody(ctx)
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

	headers = readHeaders(ctx)
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

func readBody(ctx *app.WebhookTriggersContext) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(ctx.Body)
	body := buf.String()
	return body
}

func readHeaders(ctx *app.WebhookTriggersContext) map[string][]string {
	return ctx.RequestData.Header
}
