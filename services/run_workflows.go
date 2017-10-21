package services

import (
	"encoding/json"

	"github.com/spf13/afero"

	"github.com/rchampourlier/letto_go/events"
	"github.com/rchampourlier/letto_go/exec"
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
	err = runner.Execute(event.Group, ctx)
	if err != nil {
		return err
	}

	return nil
}

func eventToExecContext(event events.ReceivedWebhook) (exec.Context, error) {
	var ctx = exec.Context{}

	// TODO: determine the type of content using the headers.ContentType
	//   and select the appropriate parser.
	parsedBody, err := parseJSONBody(event.Body)
	if err != nil {
		return ctx, err
	}

	req := exec.WebhookRequest{
		Method:  event.Method,
		URL:     event.URL,
		Host:    event.Host,
		Body:    parsedBody,
		Headers: event.Headers,
	}
	ctx = exec.Context{
		Trigger: exec.Webhook,
		Request: req,
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
