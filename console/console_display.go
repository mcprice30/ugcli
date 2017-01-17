// ugCLI is a library built atop termbox for creating CLI applications.
package console

import (
	tb "github.com/nsf/termbox-go"
)

const cursorFmt = tb.ColorDefault | tb.AttrReverse

func (c *Console) incrementCursor() {
	c.cursorX++
	c.cursorY += c.cursorX / c.width
	c.cursorX = c.cursorX % c.width
	if c.cursorY >= c.top+c.height {
		c.scrollDown()
	}
}

func (c *Console) moveCursorLeft() {
	if loc := c.getCursorLoc(); loc == 0 {
		return
	}
	tb.SetCell(c.cursorX, c.cursorY, c.getCursorChar(), tb.ColorDefault, tb.ColorDefault)
	c.decrementCursor()
	tb.SetCell(c.cursorX, c.cursorY, c.getCursorChar(), cursorFmt, cursorFmt)
}

func (c *Console) moveCursorRight() {
	if loc := c.getCursorLoc(); loc >= len(c.currline) {
		return
	}
	tb.SetCell(c.cursorX, c.cursorY, c.getCursorChar(), tb.ColorDefault, tb.ColorDefault)
	c.incrementCursor()
	tb.SetCell(c.cursorX, c.cursorY, c.getCursorChar(), cursorFmt, cursorFmt)
}

func (c *Console) getCursorLoc() int {
	return (c.cursorY-c.promptY)*c.width + c.cursorX - c.left - len(c.prompt)
}

func (c *Console) getCursorChar() rune {
	if loc := c.getCursorLoc(); loc < 0 || loc >= len(c.currline) {
		return ' '
	} else {
		return rune(c.currline[loc])
	}
}

func (c *Console) decrementCursor() {
	c.cursorX--
	if c.cursorX < 0 {
		c.cursorY--
		c.cursorX = c.left + c.width
	}
}

func (c *Console) writeChar(ch rune) {
	tb.SetCell(c.cursorX, c.cursorY, ch, tb.ColorDefault, tb.ColorDefault)
	c.incrementCursor()
	tb.SetCell(c.cursorX, c.cursorY, ' ', tb.ColorWhite, tb.ColorWhite)
}

func (c *Console) insertChar(ch rune) {
	loc := c.getCursorLoc()
	cX := c.cursorX
	cY := c.cursorY
	tb.SetCell(cX, cY, ch, tb.ColorDefault, tb.ColorDefault)
	c.incrementCursor()

	for i := loc; i < len(c.currline); i++ {
		tb.SetCell(c.cursorX, c.cursorY, rune(c.currline[i]), tb.ColorDefault, tb.ColorDefault)
		c.incrementCursor()
	}
	c.cursorX = cX
	c.cursorY = cY
	c.currline = c.currline[:loc] + string(ch) + c.currline[loc:]
	c.moveCursorRight()
}

// Print a string, with no newline, to a given Console.
func (c *Console) Print(str string) {
	for _, ch := range str {
		c.writeChar(ch)
	}
}

// Print a string, followed by a newline, to a given Console.
func (c *Console) Println(str string) {
	c.Print(str)
	tb.SetCell(c.cursorX, c.cursorY, ' ', tb.ColorDefault, tb.ColorDefault)
	c.cursorX = c.left
	c.cursorY++
	if c.cursorY >= c.top+c.height {
		c.scrollDown()
	}
	tb.SetCell(c.cursorX, c.cursorY, ' ', tb.ColorWhite, tb.ColorWhite)
}

// The equivalent of pressing the backspace key.
// Moves the cursor one cell back, then deletes the value under the cursor.
func (c *Console) backspace() {
	loc := c.getCursorLoc()
	if loc <= 0 {
		return
	}

	cX := c.cursorX
	cY := c.cursorY

	for i := loc; i+1 < len(c.currline); i++ {
		tb.SetCell(c.cursorX, c.cursorY, rune(c.currline[i+1]), tb.ColorDefault, tb.ColorDefault)
		c.incrementCursor()
	}
	tb.SetCell(c.cursorX, c.cursorY, ' ', tb.ColorDefault, tb.ColorDefault)

	if loc == 1 {
		c.currline = c.currline[loc:]
	} else if loc == len(c.currline) {
		c.currline = c.currline[:loc-1]
	} else {
		c.currline = c.currline[:loc-1] + c.currline[loc:]
	}

	c.cursorX = cX
	c.cursorY = cY
	c.moveCursorLeft()
}

// Clear all text to the left of the current line's prompt string.
func (c *Console) clearLine() {
	for _ = range c.currline {
		c.backspace()
	}
	c.currline = ""
}

// Scroll down one cell on the console.
func (c *Console) scrollDown() {
	buffer := tb.CellBuffer()
	tbWidth, _ := tb.Size()

	for y := c.top; y < c.top+c.height-1; y++ {
		for x := c.left; x < c.left+c.width; x++ {
			oldCell := buffer[(y+1)*tbWidth+x]
			tb.SetCell(x, y, oldCell.Ch, oldCell.Fg, oldCell.Bg)
		}
	}

	for x := c.left; x < c.left+c.width; x++ {
		tb.SetCell(x, c.top+c.height-1, ' ', tb.ColorDefault, tb.ColorDefault)
	}

	c.cursorY--
	c.promptY--
	if c.cursorY < 0 {
		c.cursorY = 0
	}
}
