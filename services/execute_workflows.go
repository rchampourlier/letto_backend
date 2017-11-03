package services

import (
	"gitlab.com/letto/letto_backend/exec/js"
	"gitlab.com/letto/letto_backend/services/events"
)

// ExecuteWorkflows stores the service's config.
type ExecuteWorkflows struct {
	eventBus *EventBus
	jsRunner js.Runner
}

// NewExecuteWorkflows creates a new ExecuteWorkflowsService with
// the specified `afero.Fs` filesystem abstraction.
func NewExecuteWorkflows(eventBus *EventBus, jsRunner js.Runner) *ExecuteWorkflows {
	return &ExecuteWorkflows{
		eventBus: eventBus,
		jsRunner: jsRunner,
	}
}

func (s *ExecuteWorkflows) Consume(e events.Event) error {
	typedEvent := e.(events.ActivatedTriggerEvent)
	output, err := s.jsRunner.Execute(typedEvent)
	if err != nil {
		return err
	}

	ctx := events.ExecutedWorkflowsContext{
		Stdout: output.Stdout,
		Stderr: output.Stderr,
	}
	if err != nil {
		ctx.Error = err.Error()
	}

	ed := e.Data()
	newEvent := events.NewExecutedWorkflowsEvent(ed.SequenceID, ed.Group, ctx)

	s.eventBus.Publish(newEvent)
	return nil
}

// StartConsuming
func (s *ExecuteWorkflows) StartConsuming() {
	s.eventBus.registerConsumer(s, []string{"activated_trigger"})
}

// StopConsuming
func (s *ExecuteWorkflows) StopConsuming() {
	s.eventBus.unregisterConsumer(s)
}
