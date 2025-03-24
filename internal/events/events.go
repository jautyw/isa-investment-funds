package events

type Event struct {
	broker Broker
}

// NewEvents represents a new instance of Events
func NewEvents(broker Broker) *Event {
	return &Event{
		broker: broker,
	}
}

// Broker concerns a collection of methods that concern message queues
type Broker interface {
	Publish()
}
