package controllers

import (
	"bytes"
	"time"

	"github.com/goadesign/goa"
	"github.com/satori/go.uuid"
	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/app"
	"github.com/rchampourlier/letto_go/events"
	"github.com/rchampourlier/letto_go/exec"
	"github.com/rchampourlier/letto_go/services"
	"github.com/rchampourlier/letto_go/util"
)

// TriggersController implements the triggers resource.
type TriggersController struct {
	*goa.Controller
	fs       afero.Fs
	jsRunner exec.JsRunner
}

// NewTriggersController creates a triggers controller.
func NewTriggersController(service *goa.Service, fs afero.Fs, jsRunner exec.JsRunner) *TriggersController {
	return &TriggersController{
		Controller: service.NewController("TriggersController"),
		fs:         fs,
		jsRunner:   jsRunner,
	}
}

// Webhook runs the webhook action.
func (c *TriggersController) Webhook(ctx *app.WebhookTriggersContext) error {
	// TriggersController_Webhook: start_implement

	event := events.ReceivedWebhook{
		UniqueID: uniqueID(),
		Method:   ctx.Method,
		URL:      ctx.URL,
		Host:     ctx.Host,
		Body:     readBody(ctx),
		Headers:  readHeaders(ctx),
		Group:    ctx.Group,
	}
	services.NewTrace(c.fs).OnReceivedWebhook(event)
	services.NewRunWorkflows(c.fs).OnReceivedWebhook(event)

	// TriggersController_Webhook: end_implement
	return nil
}

func uniqueID() string {
	timestamp := util.Timestamp(time.Now())
	uuid := uuid.NewV4()
	return timestamp + "-" + uuid.String()
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
