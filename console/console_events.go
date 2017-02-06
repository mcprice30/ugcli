package console

// console_events.go contains utility functions to handle termbox events within
// a console component.

import (
	"sort"

	tb "github.com/nsf/termbox-go"

	"github.com/mcprice30/ugcli"
)

const column_pad = 2

// Run will be called to launch the console. It serves as the main activity
// loop for the console, and implements the component interface, allowing
// consoles to be embedded within an ugcli application.
func (c *Console) Run(eq *ugcli.EventQueue) {

	// Print the prompt for the first time.
	c.Print(c.prompt)

	// If we don't have an executer specified, simply use an echo executer.
	if c.executer == nil {
		c.executer = DefaultExecuter(c)
	}

	// Loop until finished.
	for c.running {

		// In the event of an error, panic.
		if err := tb.Flush(); err != nil {
			panic(err)
		}

		// Get an event from the event queue.
		event := eq.PollEvent()

		// Delegate to the appropriate helper.
		if event.Type == tb.EventKey {
			switch event.Key {
			case 0:
				c.insertChar(event.Ch)
			case tb.KeySpace:
				c.insertChar(' ')
			case tb.KeyEnter:
				for i := c.getCursorLoc(); i < len(c.currline); i++ {
					c.moveCursorRight()
				}
				c.executeLine()
				c.promptY = c.cursorY
			case tb.KeyCtrlC:
				c.running = false
			case tb.KeyBackspace, tb.KeyBackspace2:
				c.backspace()
			case tb.KeyArrowUp:
				c.doArrowUp()
			case tb.KeyArrowDown:
				c.doArrowDown()
			case tb.KeyTab:
				c.doTabCompletion()
			case tb.KeyArrowRight:
				c.moveCursorRight()
			case tb.KeyArrowLeft:
				c.moveCursorLeft()
			}
		}
	}
}

// executeLine will execute the current line, go to a new line, and print a
// prompt for the next line.
func (c *Console) executeLine() {
	c.Println("")
	if c.executer != nil {
		_, c.running = c.executer.Execute(c.currline)
	}
	c.Print(c.prompt)
	if len(c.currline) > 0 {
		c.lineBuffer[c.bufferIdx%bufferSize] = c.currline
		c.bufferIdx++
	}
	c.currline = ""
	c.diff = 0
	c.oldLineCopy = ""
}

// doArrowDown will set the current line to a more recently executed command,
// or what the user was typing before pressing the up arrow, if applicable.
func (c *Console) doArrowDown() {
	if c.diff < -1 {
		c.diff++
		c.clearLine()
		c.currline = c.lineBuffer[(c.bufferIdx+c.diff)%bufferSize]
		c.Print(c.currline)
	} else if c.diff == -1 {
		c.diff++
		c.clearLine()
		c.currline = c.oldLineCopy
		c.Print(c.currline)
	}
}

// doArrowUp will set the current line to a less recently executed command.
func (c *Console) doArrowUp() {
	if bufferSize+c.diff > 0 && c.bufferIdx+c.diff > 0 {
		if c.diff == 0 {
			c.oldLineCopy = c.currline
		}
		c.diff--
		c.clearLine()
		c.currline = c.lineBuffer[(c.bufferIdx+c.diff)%bufferSize]
		c.Print(c.currline)
	}
}

// doTabCompletion will ask the user-defined completer for recommendations
// for the current line, before displaying them, if applicable.
func (c *Console) doTabCompletion() {
	if c.completer != nil {
		prefix, options := c.completer.Complete(c.currline)

		if len(options) > 1 {
			c.clearLine()
			c.currline = prefix
			c.Print(c.currline)
			c.Println("")
			c.printOptions(options)
			c.promptY = c.cursorY
			c.Print(c.prompt)
			c.Print(c.currline)
		} else if len(options) == 1 {
			c.clearLine()
			c.currline = prefix
			c.Print(c.currline)
		}
	}
}

// printOptions is a utility function that will print a variety of strings
// in columns.
func (c *Console) printOptions(options []string) {
	sort.Strings(options)
	maxLen := 0
	for _, option := range options {
		if l := len(option); maxLen < l {
			maxLen = l
		}
	}
	numColumns := (c.width + column_pad) / (maxLen + column_pad)
	if numColumns == 0 {
		numColumns++
	}
	printed := 0
	for _, option := range options {
		c.Print(option)
		for i := len(option); i < maxLen+column_pad; i++ {
			c.Print(" ")
		}
		printed++
		if printed%numColumns == 0 {
			c.Println("")
		}
	}
	if printed%numColumns != 0 {
		c.Println("")
	}
}
