package services

import (
	events "gitlab.com/letto/letto_backend/services/events"
)

// Service is the interface that must be fulfilled by the `Service`
// struct.
type Service interface {
	Consumer
	StartConsuming()
	StopConsuming()
}

type Consumer interface {
	Consume(e events.Event) error
}
