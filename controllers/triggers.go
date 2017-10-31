package controllers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/goadesign/goa"
	"github.com/satori/go.uuid"
	"github.com/spf13/afero"

	"gitlab.com/letto/letto_backend/app"
	"gitlab.com/letto/letto_backend/events"
	"gitlab.com/letto/letto_backend/exec/js"
	"gitlab.com/letto/letto_backend/services"
	"gitlab.com/letto/letto_backend/util"
)

// TriggersController implements the triggers resource.
type TriggersController struct {
	*goa.Controller
	fs       afero.Fs
	jsRunner js.Runner
}

// NewTriggersController creates a triggers controller.
func NewTriggersController(service *goa.Service, fs afero.Fs, jsRunner js.Runner) *TriggersController {
	return &TriggersController{
		Controller: service.NewController("TriggersController"),
		fs:         fs,
		jsRunner:   jsRunner,
	}
}

// Webhook runs the webhook action.
func (c *TriggersController) Webhook(ctx *app.WebhookTriggersContext) error {
	// TriggersController_Webhook: start_implement
	var err error

	event := events.ReceivedWebhook{
		UniqueID: uniqueID(),
		Method:   ctx.Method,
		URL:      ctx.URL,
		Host:     ctx.Host,
		Body:     readBody(ctx),
		Headers:  readHeaders(ctx),
		Group:    ctx.Group,
	}

	// TODO: improve error management, should not have to print them here
	//   but since we only return once not to fail the webhook's call,
	//   it's not ideal.
	err = services.NewTrace(c.fs).OnReceivedWebhook(event)
	if err != nil {
		fmt.Printf("Trace.OnReceivedWebhook error: %s\n", err)
	}
	err = services.NewRunWorkflows(c.fs, c.jsRunner).OnReceivedWebhook(event)
	if err != nil {
		fmt.Printf("RunWorkflows.OnReceivedWebhook error: %s\n", err)
	}

	// TriggersController_Webhook: end_implement
	return err
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
