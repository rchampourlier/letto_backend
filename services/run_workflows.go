package services

import (
	"encoding/json"

	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_backend/events"
	"github.com/rchampourlier/letto_backend/exec"
	"github.com/rchampourlier/letto_backend/exec/values"
)

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
	ctx, err := eventToExecContext(event)
	if err != nil {
		return err
	}

	runner := exec.NewJsRunner(s.Fs)
	output, err := runner.Execute(event.Group, ctx)

	newEvent := events.CompletedWorkflows{
		TriggerUniqueID: event.UniqueID,
		Group:           event.Group,
		Stdout:          output.Stdout,
		Stderr:          output.Stderr,
	}
	if err != nil {
		newEvent.Error = err.Error()
	}
	NewTrace(s.Fs).OnCompletedWorkflows(newEvent)

	return nil
}

func eventToExecContext(event events.ReceivedWebhook) (values.Context, error) {
	var ctx = values.Context{}

	// TODO: determine the type of content using the headers.ContentType
	//   and select the appropriate parser.
	parsedBody, err := parseJSONBody(event.Body)
	if err != nil {
		return ctx, err
	}

	req := values.WebhookRequest{
		Method:  event.Method,
		URL:     event.URL,
		Host:    event.Host,
		Body:    parsedBody,
		Headers: event.Headers,
	}
	ctx = values.Context{
		Trigger: values.Webhook,
		Request: req,
		Group:   event.Group,
	}
	return ctx, nil
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
