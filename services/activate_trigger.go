package services

import (
	"log"

	"gitlab.com/letto/letto_backend/services/events"
)

// ActivateTrigger is a service consuming `ReceivedWebhook` and
// `ActivatedSchedule` events and producing `ActivatedTrigger`
// events when applicable.
type ActivateTrigger struct {
	eventBus *EventBus
}

// NewActivateTrigger creates a new ExecuteWorkflowsService with
// the specified `afero.Fs` filesystem abstraction.
func NewActivateTrigger(eventBus *EventBus) *ActivateTrigger {
	return &ActivateTrigger{
		eventBus: eventBus,
	}
}

// Consume is called by the event-bus when a consumed event is
// received.
func (s *ActivateTrigger) Consume(e events.Event) error {
	var newEvent events.Event
	ed := e.Data()
	switch e.(type) {
	case events.ReceivedWebhookEvent:
		newEvent = events.NewActivatedTriggerEvent(
			ed.ID,
			ed.Group,
			events.ActivatedTriggerContext{e},
		)
	default:
		log.Fatalf("Unknown event type %s\n", e)
	}
	s.eventBus.Publish(newEvent)
	return nil
}

// StartConsuming
func (s *ActivateTrigger) StartConsuming() {
	s.eventBus.registerConsumer(s, []string{"received_webhook"})
}

// StopConsuming
func (s *ActivateTrigger) StopConsuming() {
	s.eventBus.unregisterConsumer(s)
}
