package controllers

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/goadesign/goa"

	"gitlab.com/letto/letto_backend/app"
	"gitlab.com/letto/letto_backend/exec/js"
	"gitlab.com/letto/letto_backend/services"
	"gitlab.com/letto/letto_backend/services/events"
)

// TriggersController implements the triggers resource.
type TriggersController struct {
	*goa.Controller
	eventBus *services.EventBus
	jsRunner js.Runner
}

// NewTriggersController creates a triggers controller.
func NewTriggersController(service *goa.Service, eventBus *services.EventBus, jsRunner js.Runner) *TriggersController {
	return &TriggersController{
		Controller: service.NewController("TriggersController"),
		eventBus:   eventBus,
		jsRunner:   jsRunner,
	}
}

// Webhook runs the webhook action.
func (c *TriggersController) Webhook(ctx *app.WebhookTriggersContext) error {
	// TriggersController_Webhook: start_implement
	var err error

	body := readBody(ctx)

	// TODO: should support non-JSON bodies
	parsedBody, err := parseJSONBody(body)
	if err != nil {
		return err
	}

	eventCtx := events.ReceivedWebhookContext{
		ReceptionTime: time.Now(),
		Method:        ctx.Method,
		URL:           ctx.URL,
		Host:          ctx.Host,
		Body:          parsedBody,
		Headers:       readHeaders(ctx),
	}
	eventData := events.NewEventData(ctx.Group, eventCtx)
	event := events.ReceivedWebhookEvent{
		EventData: eventData,
	}
	c.eventBus.Publish(event)

	// TriggersController_Webhook: end_implement
	return err
}

func readBody(ctx *app.WebhookTriggersContext) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(ctx.Body)
	body := buf.String()
	return body
}

func parseJSONBody(body string) (map[string]interface{}, error) {
	var bodyParsed map[string]interface{}
	if len(body) == 0 {
		return nil, nil
	}

	err := json.Unmarshal([]byte(body), &bodyParsed)
	if err != nil {
		return nil, err
	}
	return bodyParsed, nil
}

func readHeaders(ctx *app.WebhookTriggersContext) map[string][]string {
	return ctx.RequestData.Header
}
