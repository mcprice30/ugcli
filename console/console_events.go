// ugCLI is a library built atop termbox for creating CLI applications.
package console

import (
	"sort"

	tb "github.com/nsf/termbox-go"

	"github.com/mcprice30/ugcli"
)

const column_pad = 2

// Call this to launch the console. Main activity loop for
// the console.
func (c *Console) Run(eq *ugcli.EventQueue) {

	c.promptY = c.cursorY
	c.Print(c.prompt)

	if c.executer == nil {
		c.executer = DefaultExecuter(c)
	}

	for c.running {

		if err := tb.Flush(); err != nil {
			panic(err)
		}

		event := eq.PollEvent()
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

func (c *Console) executeLine() {
	c.Println("")
	if c.executer != nil {
		_, c.running = c.executer.Execute(c.currline)
	}
	c.Print(c.prompt)
	if len(c.currline) > 0 {
		c.lineBuffer[c.bufferIdx%buffer_size] = c.currline
		c.bufferIdx++
	}
	c.currline = ""
	c.diff = 0
	c.oldLineCopy = ""
}

func (c *Console) doArrowDown() {
	if c.diff < -1 {
		c.diff++
		c.clearLine()
		c.currline = c.lineBuffer[(c.bufferIdx+c.diff)%buffer_size]
		c.Print(c.currline)
	} else if c.diff == -1 {
		c.diff++
		c.clearLine()
		c.currline = c.oldLineCopy
		c.Print(c.currline)
	}
}

func (c *Console) doArrowUp() {
	if buffer_size+c.diff > 0 && c.bufferIdx+c.diff > 0 {
		if c.diff == 0 {
			c.oldLineCopy = c.currline
		}
		c.diff--
		c.clearLine()
		c.currline = c.lineBuffer[(c.bufferIdx+c.diff)%buffer_size]
		c.Print(c.currline)
	}
}

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
