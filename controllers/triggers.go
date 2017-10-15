package controllers

import (
	"bytes"

	"github.com/goadesign/goa"
	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/app"
	"github.com/rchampourlier/letto_go/events"
	"github.com/rchampourlier/letto_go/services"
)

// TriggersController implements the triggers resource.
type TriggersController struct {
	*goa.Controller
	fs afero.Fs
}

// NewTriggersController creates a triggers controller.
func NewTriggersController(service *goa.Service, fs afero.Fs) *TriggersController {
	return &TriggersController{
		Controller: service.NewController("TriggersController"),
		fs:         fs,
	}
}

// Webhook runs the webhook action.
func (c *TriggersController) Webhook(ctx *app.WebhookTriggersContext) error {
	// TriggersController_Webhook: start_implement

	event := events.ReceivedWebhook{
		Method:  ctx.Method,
		URL:     ctx.URL,
		Host:    ctx.Host,
		Body:    readBody(ctx),
		Headers: readHeaders(ctx),
		Group:   ctx.Group,
	}
	services.NewTrace(c.fs).OnReceivedWebhook(event)
	services.NewRunWorkflows(c.fs).OnReceivedWebhook(event)

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
