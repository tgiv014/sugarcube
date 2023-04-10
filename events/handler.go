package events

import (
	"io"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type eventContainer struct {
	Topic   string
	Payload Event
}

func (b *Bus) Handler(c *gin.Context) {
	eventChan, unsubscribe := b.Subscribe("*")
	defer unsubscribe()

	// Emit an initial connection message
	go func() {
		err := b.Emit(BusEvent{"Connected"})
		log.Info("Sent connection message", "err", err)
	}()

	log.Info("Opening SSE event bus")
	defer log.Info("SSE event bus shutdown")
	c.Stream(func(w io.Writer) bool {
		if ev, ok := <-eventChan; ok {
			c.SSEvent("message", eventContainer{ev.Topic(), ev})
			return true
		}
		return false
	})
}
