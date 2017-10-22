package values

import (
	"net/url"
)

// TriggerType represents the type of trigger at the origin
// of the execution. May be webhook only for now (scheduler
// will be added in the future).
type TriggerType int

const (
	// Webhook when the execution has been triggered by a
	// webhook received.
	Webhook TriggerType = iota
)

// Context represents the context of a code execution
// perform by an `exec` module.
type Context struct {
	Trigger TriggerType
	Request WebhookRequest
	Group   string
}

// WebhookRequest represents the data of a webhook request.
type WebhookRequest struct {
	Method  string
	URL     *url.URL
	Host    string
	Body    map[string]interface{}
	Headers map[string][]string
}
