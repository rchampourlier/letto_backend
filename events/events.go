package events

import (
	"net/url"
)

// ReceivedWebhook is the data struct for the event
// fired on receiving a webhook.
type ReceivedWebhook struct {
	Method  string
	URL     *url.URL
	Host    string
	Body    string
	Headers map[string][]string
	Group   string
}
