// package ugcli contains interfaces for creating ugly command line interface
// applications.
package ugcli

import (
	tb "github.com/nsf/termbox-go"
)

// Cli represents a CLI application. It contains a variety of sub-components
// that can be combined to form a larger application.
type Cli struct {

	// activeComponent stores the index of which component/handler pair is
	// currently active.
	activeComponent int

	// components stores all sub-components that comprise this application.
	components []Component

	// handlers stores the event queues to delegate events to various components.
	handlers []*EventQueue

	// runningComponents indicates the number of components currently running.
	runningComponents int

	// eventBuffer is a channel that will grab events from the termbox event poll.
	eventBuffer chan tb.Event

	// doneChan will send kill signals to the application.
	doneChan chan bool
}

// NewCli will create a new CLI application.
func NewCli() *Cli {
	return &Cli{
		activeComponent:   -1,
		components:        []Component{},
		handlers:          []*EventQueue{},
		runningComponents: 0,
		eventBuffer:       make(chan tb.Event, 10),
		doneChan:          make(chan bool),
	}
}

// AddComponent will bind a subcomponent to this application.
func (c *Cli) AddComponent(comp Component) {
	c.components = append(c.components, comp)
	c.handlers = append(c.handlers, newEventQueue())
	c.activeComponent = 0
}

// eventPoll serves as a background goroutine to listen for events from termbox.
func (c *Cli) eventPoll() {
	for {
		c.eventBuffer <- tb.PollEvent()
	}
}

// Run launches the ugcli application.
func (c *Cli) Run() {
	for i := range c.components {
		c.runningComponents++
		// Launch each component in a new thread.
		go c.runComponent(i)
	}

	// Start listening for termbox events.
	go c.eventPoll()

	for {

		// Wait for one of the channels to get an event.
		select {

		// Either it comes from termbox, in which case it must be delegated to
		// the appropriate component's event buffer.
		case event := <-c.eventBuffer:
			c.handlers[c.activeComponent].addEvent(event)
		// Or it is a kill signal, in which case we should exit the application.
		case done := <-c.doneChan:
			if done {
				return
			}
		}
	}
}

// runComponent is a wrapper thread for a given component.
func (c *Cli) runComponent(comp int) {
	c.components[comp].Run(c.handlers[comp])

	c.runningComponents--
	if c.runningComponents == 0 {
		// Once the component stops running, send a kill signal iff this was the
		// last running component.
		c.doneChan <- true
	}
}

// Component represents a subcomponent that can be added to the ugcli app.
type Component interface {
	Run(*EventQueue)
}
