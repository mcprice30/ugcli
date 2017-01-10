// ugCLI is a library built atop termbox for creating CLI applications.
package ugcli

import (
	tb "github.com/nsf/termbox-go"
)

func (c *Console) writeChar(ch rune) {
	tb.SetCell(c.cursorX, c.cursorY, ch, tb.ColorDefault, tb.ColorDefault)
	c.cursorX++
	c.cursorY += c.cursorX / c.width
	c.cursorX = c.cursorX % c.width
	if c.cursorY >= c.top+c.height {
		c.scrollDown()
	}
	tb.SetCell(c.cursorX, c.cursorY, ' ', tb.ColorWhite, tb.ColorWhite)
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
	tb.SetCell(c.cursorX, c.cursorY, ' ', tb.ColorDefault, tb.ColorDefault)
	c.cursorX--
	if c.cursorX < 0 {
		c.cursorY--
		c.cursorX = c.left + c.width
	}
	tb.SetCell(c.cursorX, c.cursorY, ' ', tb.ColorWhite, tb.ColorWhite)
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
	if c.cursorY < 0 {
		c.cursorY = 0
	}
}
