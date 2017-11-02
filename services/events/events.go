package events

import (
	"net/url"
	"time"

	uuid "github.com/goadesign/goa/uuid"
)

type Event interface {
	Name() string
	Data() EventData
}

// EventData contains the common values for events.
type EventData struct {
	ID      uuid.UUID
	Time    time.Time
	Group   string
	Context interface{}
}

type ReceivedWebhookEvent struct {
	EventData                // embedded struct
	receivedWebhookContexter // embedded interface
}

// ReceivedWebhookContext contains the context of a
// `ReceivedWebhook` event.
type ReceivedWebhookContext struct {
	ReceptionTime time.Time
	Method        string
	URL           *url.URL
	Host          string
	Body          map[string]interface{}
	Headers       map[string][]string
}

func (e ReceivedWebhookEvent) ReceivedWebhookContext() ReceivedWebhookContext {
	return e.Context.(ReceivedWebhookContext)
}

func (e ReceivedWebhookEvent) Name() string {
	return "received_webhook"
}

type ActivatedScheduleEvent struct {
	EventData
	activatedScheduleContexter
}

func (e ActivatedScheduleEvent) ActivatedScheduleContext() ActivatedScheduleContext {
	return e.Context.(ActivatedScheduleContext)
}

func (e ActivatedScheduleEvent) Name() string {
	return "activated_schedule"
}

type ActivatedTriggerEvent struct {
	EventData
	activatedTriggerContexter
}

func (e ActivatedTriggerEvent) ActivatedTriggerContext() ActivatedTriggerContext {
	return e.Context.(ActivatedTriggerContext)
}

func (e ActivatedTriggerEvent) Name() string {
	return "activated_trigger"
}

type ExecutedWorkflowsEvent struct {
	EventData
	executedWorkflowsContexter
}

// ExecutedWorkflowsContext is the data struct for the event
// fired when the workflows have been completed.
type ExecutedWorkflowsContext struct {
	Stdout string
	Stderr string
	Error  string
}

func (e ExecutedWorkflowsEvent) ExecutedWorkflowsContext() ExecutedWorkflowsContext {
	return e.Context.(ExecutedWorkflowsContext)
}

func (e ExecutedWorkflowsEvent) Name() string {
	return "executed_workflows"
}

func (ed EventData) Data() EventData {
	return ed
}

// NewEventData returns a new `EventData` struct with a randomly
// generate ID (an UUID) and the current time, as well as the
// specified group and context.
func NewEventData(group string, ctx interface{}) EventData {
	return EventData{
		ID:      uuid.NewV4(),
		Time:    time.Now(),
		Group:   group,
		Context: ctx,
	}
}

// ActivatedScheduleContext is the event fired when a schedule is being
// activated become the time has come.
type ActivatedScheduleContext struct {
	ActivationTime time.Time
	ScheduleID     uuid.UUID
	ScheduleData   map[string]interface{}
}

// ActivatedTriggerContext is the event fired when a trigger has been
// activated. A trigger is activated when a webhook has been
// received and can activate a trigger or when a schedule is
// activated.
type ActivatedTriggerContext struct {
	ActivatingEvent Event
}

type receivedWebhookContexter interface {
	ReceivedWebhookContext() ReceivedWebhookContext
}
type activatedScheduleContexter interface {
	ActivatedScheduleContext() ActivatedScheduleContext
}

type activatedTriggerContexter interface {
	ActivatedTriggerContext() ActivatedTriggerContext
}

type executedWorkflowsContexter interface {
	ExecutedWorkflowsContext() ExecutedWorkflowsContext
}
