package events

import "errors"

// Event is any object that can provide a topic, encode itself a json
type Event interface {
	Topic() string
	Decode(v any) error
}

type BusEvent struct {
	Message string
}

func (b BusEvent) Topic() string {
	return "busEvent"
}

func (b BusEvent) Decode(v any) error {
	return errors.New("not implemented")
}
