// ugCLI is a library built atop termbox for creating CLI applications.
package ugcli

import (
	tb "github.com/nsf/termbox-go"
)

const buffer_size = 100

// Call this to launch the console. Main activity loop for
// the console.
func (c *Console) Run() {
	running := true

	c.Print(c.prompt)

	currline := ""
	lineBuffer := make([]string, buffer_size)
	bufferIdx := 0
	diff := 0

	oldLineCopy := ""

	for running {

		if err := tb.Flush(); err != nil {
			panic(err)
		}

		event := tb.PollEvent()
		if event.Type == tb.EventKey {
			switch event.Key {
				case 0:
					c.writeChar(event.Ch)
					currline += string(event.Ch)
				case tb.KeySpace:
					c.writeChar(' ')
					currline += " "
				case tb.KeyEnter:
					c.Println("")
					running = c.executeLine(currline)
					c.Print(c.prompt)
					lineBuffer[bufferIdx % buffer_size] = currline
					bufferIdx++
					currline = ""
					diff = 0
					oldLineCopy = ""
				case tb.KeyCtrlC:
					running = false
				case tb.KeyBackspace, tb.KeyBackspace2:
					if len(currline) > 0 {
						c.backspace()
						currline = currline[:len(currline)-1]
					}
				case tb.KeyArrowUp:
					if buffer_size + diff > 0 && bufferIdx + diff > 0 {
						if diff == 0 {
							oldLineCopy = currline
						}
						diff--
						for _ = range currline {
							c.backspace()
						}
						currline = lineBuffer[(bufferIdx + diff) % buffer_size]
						c.Print(currline)
				  }
				case tb.KeyArrowDown:
					if diff < -1 {
						diff++
						for _ = range currline {
							c.backspace()
						}
						currline = lineBuffer[(bufferIdx + diff) % buffer_size]
						c.Print(currline)
					} else if diff == -1 {
						diff++
						for _ = range currline {
							c.backspace()
						}
						currline = oldLineCopy
						c.Print(currline)
					}
			}
		}
	}
}
