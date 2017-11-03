// Warning: this package is not simple and a bit verbose.
//
// The `...Contexter` interfaces have been created to enable
// event structs (like `ReceivedWebhookEvent`) to have different
// interfaces so they could be distinguished in a type case. This
// may not be really necessary.
// TODO: try to remove the `...Contexter` interfaces.
//
// The use of an `interface` for `Event` clearly enables the
// `switch event.(type)` structure that is being used in some
// services consuming several events but needing to process
// them differently (e.g. `ActivateTrigger`).
//
// Each event requires several structures to be written, making
// the whole thing quite verbose. A code-generation approach
// like Goa may be a solution.
//
// The `EventData` struct is the only struct containing event
// data, since `Event` is only an interface. Other event types,
// like `ReceivedWebhookEvent`, must embed an `EventData`.
// This allows all modules dealing with events to manipulate the
// same structure. The custom part is contained in `Context` which
// is an `interface{}` struct.
//
// The `...Context()` method provided by each event (e.g.
// `ReceivedWebhookContext()` enables to access the event's
// context with the appropriate type (e.g. `ReceivedWebhookContext`).
// Each event type must implement this method for its own
// context type.
package events

import (
	"fmt"
	"log"
	"net/url"
	"time"

	uuid "github.com/goadesign/goa/uuid"
	"gitlab.com/letto/letto_backend/util"
)

// Event is the interface for event structures.
type Event interface {
	FullIdentifier() string
	Data() EventData
}

// EventData contains the common values for events.
type EventData struct {
	Name       string
	SequenceID uuid.UUID
	ID         uuid.UUID
	Time       time.Time
	Group      string
	Context    interface{}
}

// ReceivedWebhookEvent is the event fired when a webhook
// as been received by the `triggers/webhook` endpoint.
type ReceivedWebhookEvent struct {
	EventData                // embedded struct
	receivedWebhookContexter // embedded interface
}

func NewReceivedWebhookEvent(group string, ctx interface{}) ReceivedWebhookEvent {
	return ReceivedWebhookEvent{
		EventData: newEventData("received_webhook", nil, group, ctx),
	}
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

type receivedWebhookContexter interface {
	ReceivedWebhookContext() ReceivedWebhookContext
}

type ActivatedScheduleEvent struct {
	EventData
	activatedScheduleContexter
}

// ActivatedScheduleContext is the event fired when a schedule is being
// activated become the time has come.
type ActivatedScheduleContext struct {
	ActivationTime time.Time
	ScheduleID     uuid.UUID
	ScheduleData   map[string]interface{}
}

func NewActivatedScheduleEvent(seqID uuid.UUID, group string, ctx interface{}) ActivatedScheduleEvent {
	return ActivatedScheduleEvent{
		EventData: newEventData("activated_schedule", []uuid.UUID{seqID}, group, ctx),
	}
}

func (e ActivatedScheduleEvent) ActivatedScheduleContext() ActivatedScheduleContext {
	return e.Context.(ActivatedScheduleContext)
}

type activatedScheduleContexter interface {
	ActivatedScheduleContext() ActivatedScheduleContext
}

func (e ActivatedScheduleEvent) Name() EventData {
	d := e.EventData
	d.Name = "activated_schedule"
	return d
}

type ActivatedTriggerEvent struct {
	EventData
	activatedTriggerContexter
}

// ActivatedTriggerContext is the event fired when a trigger has been
// activated. A trigger is activated when a webhook has been
// received and can activate a trigger or when a schedule is
// activated.
type ActivatedTriggerContext struct {
	ActivatingEvent Event
}

func NewActivatedTriggerEvent(seqID uuid.UUID, group string, ctx interface{}) ActivatedTriggerEvent {
	return ActivatedTriggerEvent{
		EventData: newEventData("activated_trigger", []uuid.UUID{seqID}, group, ctx),
	}
}

func (e ActivatedTriggerEvent) ActivatedTriggerContext() ActivatedTriggerContext {
	return e.Context.(ActivatedTriggerContext)
}

type activatedTriggerContexter interface {
	ActivatedTriggerContext() ActivatedTriggerContext
}

func (e ActivatedTriggerEvent) Name() EventData {
	d := e.EventData
	d.Name = "activated_trigger"
	return d
}

type ExecutedWorkflowsEvent struct {
	EventData
	executedWorkflowsContexter
}

func NewExecutedWorkflowsEvent(seqID uuid.UUID, group string, ctx interface{}) ExecutedWorkflowsEvent {
	return ExecutedWorkflowsEvent{
		EventData: newEventData("executed_workflows", []uuid.UUID{seqID}, group, ctx),
	}
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

func (e ExecutedWorkflowsEvent) Name() EventData {
	d := e.EventData
	d.Name = "executed_workflows"
	return d
}

func (ed EventData) Data() EventData {
	return ed
}

func (ed EventData) FullIdentifier() string {
	timestamp := util.Timestamp(ed.Time)
	return fmt.Sprintf("%s--%s--%s--%s.json", timestamp, ed.ID, ed.Name, ed.Group)
}

type executedWorkflowsContexter interface {
	ExecutedWorkflowsContext() ExecutedWorkflowsContext
}

// newEventData returns a new `EventData` struct with a randomly
// generate ID (an UUID) and the current time, as well as the
// specified group and context.
//
// `seqID` is a slice of `uuid.UUID` that may:
//   - be nil or an empty slice if the event is the 1st of the sequence, in which case the
//     `eventData`'s `SequenceID` will take the same value as the event's
//     `ID`.
//   - contain a single element.
//
// NB: if `seqID` contains more than 1 element, the program will fail and exit.
func newEventData(name string, seqID []uuid.UUID, group string, ctx interface{}) EventData {
	id := uuid.NewV4()
	var finalSeqID uuid.UUID
	if seqID == nil || len(seqID) == 0 {
		finalSeqID = id
	} else if len(seqID) > 1 {
		log.Panicf("events.newEventData(...) received a `seqID` with more than 1 item (%s)\n", seqID)
	} else {
		finalSeqID = seqID[0]
	}
	return EventData{
		Name:       name,
		SequenceID: finalSeqID, // the ID of the event initiating the seq
		ID:         id,
		Time:       time.Now(),
		Group:      group,
		Context:    ctx,
	}
}
