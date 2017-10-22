package events

import (
	"net/url"
)

// ReceivedWebhook is the data struct for the event
// fired on receiving a webhook.
type ReceivedWebhook struct {
	UniqueID string // <timestamp>-<uuid>
	Method   string
	URL      *url.URL
	Host     string
	Body     string
	Headers  map[string][]string
	Group    string
}

// CompletedWorkflows is the data struct for the event
// fired when the workflows have been completed.
type CompletedWorkflows struct {
	TriggerUniqueID string
	Group           string
	Stdout          string
	Stderr          string
	Error           string
}
