package events

import (
	"sync"

	"github.com/charmbracelet/log"
)

type Bus struct {
	lock        sync.Mutex
	counter     int
	subscribers map[string]map[int]chan Event
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[string]map[int]chan Event),
	}
}

func (b *Bus) Subscribe(topic string) (chan Event, func()) {
	b.lock.Lock()
	defer b.lock.Unlock()

	callbackIdx := b.counter
	b.counter++

	if _, ok := b.subscribers[topic]; !ok {
		b.subscribers[topic] = make(map[int]chan Event)
	}

	returnChan := make(chan Event)

	b.subscribers[topic][callbackIdx] = returnChan

	return returnChan, func() {
		b.lock.Lock()
		defer b.lock.Unlock()
		delete(b.subscribers[topic], callbackIdx)
		close(returnChan)
	}
}

func (b *Bus) Emit(event Event) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	totalSubscribers := 0

	channels, ok := b.subscribers["*"]
	if ok {
		for _, channel := range channels {
			totalSubscribers++
			channel <- event
		}
	}

	channels, ok = b.subscribers[event.Topic()]
	if ok {
		for _, channel := range channels {
			totalSubscribers++
			channel <- event
		}
	}

	log.Infof("Event %s emitted to %d listeners", event.Topic(), totalSubscribers)

	return nil
}
