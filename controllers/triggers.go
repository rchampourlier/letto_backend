package controllers

import (
	"bytes"

	"github.com/goadesign/goa"
	"github.com/rchampourlier/letto_go/app"
	"github.com/rchampourlier/letto_go/exec"
)

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

	body := readBody(ctx)
	headers := readHeaders(ctx)

	script := `
	var xmlHttp = new XMLHttpRequest();
	xmlHttp.open( "GET", "www.google.fr", false );
        xmlHttp.send( null );
        console.log(xmlHttp.responseText);
	`
	exec.RunJS(script, body, headers)

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
