package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

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
func (c *TriggersController) Webhook(ctx *app.WebhookTriggersContext) error {
	// TriggersController_Webhook: start_implement

	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	//group := ctx.Group
	//body := readBody(ctx)
	//headers := readHeaders(ctx)
	tmpDir, err := ioutil.TempDir("", tmpDirPrefix)
	if err != nil {
		return err
	}
	copyFile(path.Join(rootDir, "exec", "js", "data.js"), path.Join(tmpDir, "data.js"))
	copyFile(path.Join(rootDir, "exec", "js", "test.js"), path.Join(tmpDir, "test.js"))

	exec.RunJS(tmpDir)

	err = appFs.RemoveAll(tmpDir)
	if err != nil {
		fmt.Printf("Could not remove tmp dir: %s\n", tmpDir)
	}

	// TriggersController_Webhook: end_implement
	return nil
}

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
