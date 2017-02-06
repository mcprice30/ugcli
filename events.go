package ugcli

import (
	tb "github.com/nsf/termbox-go"
)

// EventQueue is used to pass events from termbox to individual ugcli components
type EventQueue struct {
	eventBuffer chan tb.Event
}

// newEventQueue will create a new EventQueue object.
func newEventQueue() *EventQueue {
	return &EventQueue{
		eventBuffer: make(chan tb.Event, 10),
	}
}

// send an event to the queue.
func (q *EventQueue) addEvent(e tb.Event) {
	q.eventBuffer <- e
}

// PollEvent will block until a new event is added to the queue, at which point
// it will pass it to the appropriate component.
func (q *EventQueue) PollEvent() tb.Event {
	return <-q.eventBuffer
}
