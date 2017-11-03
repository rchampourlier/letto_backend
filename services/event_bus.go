package services

import (
	"encoding/json"
	"github.com/spf13/afero"
	"log"
	"path"

	"gitlab.com/letto/letto_backend/services/events"
)

// EventBus stores an event-bus reference.
// The event-bus stores the references to registered consumers
// and the events they are consuming.
type EventBus struct {
	eventNamesConsumersMap map[string][]Consumer
	fs                     afero.Fs
	logsDir                string
}

// NewEventBus returns a new bus
func NewEventBus(fs afero.Fs, logsDir string) EventBus {
	b := EventBus{
		eventNamesConsumersMap: make(map[string][]Consumer),
		fs:      fs,
		logsDir: logsDir,
	}
	return b
}

// Publish publishes the specified event
func (eb *EventBus) Publish(e events.Event) {
	err := eb.writeEventLog(e)
	if err != nil {
		log.Printf("[ERROR] failed to log event `%s`\n", e.FullIdentifier())
	}
	log.Printf("[INFO] publishing event `%s`\n", e.FullIdentifier())
	eb.sendToConsumingServices(e)
}

// registerService register a consumer which will be called when
// the specified `listenedEvents` are published.
func (eb *EventBus) registerConsumer(c Consumer, eventNames []string) {
	for i := range eventNames {
		eventName := eventNames[i]

		eb.eventNamesConsumersMap[eventName] = append(eb.eventNamesConsumersMap[eventName], c)
	}
}

// unregisterService unregisters the specified consumer. The consumer
// will not be notified of any published events anymore.
func (eb *EventBus) unregisterConsumer(removedConsumer Consumer) {
	for event, consumers := range eb.eventNamesConsumersMap {
		remainingConsumers := make([]Consumer, 0)
		for i := range consumers {
			consumer := consumers[i]
			if consumer != removedConsumer {
				remainingConsumers = append(remainingConsumers, consumer)
			}
		}
		eb.eventNamesConsumersMap[event] = remainingConsumers
	}
}

func (eb *EventBus) writeEventLog(e events.Event) error {
	dirPath := eb.logsDir
	err := eb.fs.MkdirAll(dirPath, 0777)
	if err != nil {
		return logTraceError(err)
	}

	filePath := path.Join(dirPath, e.FullIdentifier())

	// Write the content of the event to a file
	eventAsJSON, err := json.Marshal(e.Data())
	if err != nil {
		return logTraceError(err)
	}
	err = afero.WriteFile(eb.fs, filePath, eventAsJSON, 0777)
	if err != nil {
		return logTraceError(err)
	}

	return nil
}

func (eb *EventBus) sendToConsumingServices(e events.Event) {
	consumingServices := eb.eventNamesConsumersMap[e.Data().Name]
	for i := range consumingServices {
		service := consumingServices[i]
		service.Consume(e)
	}
}

func logTraceError(err error) error {
	log.Printf("[ERROR] failed to write file `%s`\n", err)
	return err
}
