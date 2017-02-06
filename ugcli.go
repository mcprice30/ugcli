package ugcli

import (
	tb "github.com/nsf/termbox-go"
)

type Cli struct {
	activeComponent   int
	components        []Component
	handlers          []*EventQueue
	runningComponents int
	eventBuffer       chan tb.Event
	doneChan          chan bool
}

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

func (c *Cli) AddComponent(comp Component) {
	c.components = append(c.components, comp)
	c.handlers = append(c.handlers, NewEventQueue())
	c.activeComponent = 0
}

func (c *Cli) eventPoll() {
	for {
		c.eventBuffer <- tb.PollEvent()
	}
}

func (c *Cli) Run() {
	for i := range c.components {
		c.runningComponents++
		go c.runComponent(i)
	}

	go c.eventPoll()

	for {
		select {
		case event := <-c.eventBuffer:
			c.handlers[c.activeComponent].addEvent(event)
		case done := <-c.doneChan:
			if done {
				return
			}
		}
	}
}

func (c *Cli) runComponent(comp int) {
	c.components[comp].Run(c.handlers[comp])
	c.runningComponents--
	if c.runningComponents == 0 {
		c.doneChan <- true
	}
}

type Component interface {
	Run(*EventQueue)
}
