package events

type Event struct {
	broker Broker
}

func NewStore(broker Broker) *Event {
	return &Event{
		broker: broker,
	}
}

type Broker interface {
	Publish()
}
