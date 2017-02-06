package ugcli

import (
	tb "github.com/nsf/termbox-go"
)

type EventQueue struct {
	eventBuffer chan tb.Event
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		eventBuffer: make(chan tb.Event, 10),
	}
}

func (q *EventQueue) addEvent(e tb.Event) {
	q.eventBuffer <- e
}

func (q *EventQueue) PollEvent() tb.Event {
	return <-q.eventBuffer
}
