package controllers

import (
	"bytes"
	"io/ioutil"

	"github.com/goadesign/goa"
	"github.com/rchampourlier/letto_go/app"
	"github.com/rchampourlier/letto_go/exec"
)

var tmpDirPrefix = "letto"

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

	group := ctx.Group
	//body := readBody(ctx)
	//headers := readHeaders(ctx)
	tmpDir, err := ioutil.TempDir("", tmpDirPrefix)
	if err != nil {
		panic("Could not create temp directory. Abandoning.")
	}

	exec.RunJS(group, tmpDir)

	// TriggersController_Webhook: end_implement
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
