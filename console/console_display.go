package console

// console_display.go contains utility functions for actually rendering
// console elements into the terminal.

import (
	tb "github.com/nsf/termbox-go"
)

// cursorFmt is used in setting the color of the cursor itself or any text that
// the cursor is hovering over.
const cursorFmt = tb.ColorDefault | tb.AttrReverse

// incrementCursor will move the cursor one cell to the right, scrolling the
// screen if necessary, but without redrawing the cursor.
func (c *Console) incrementCursor() {
	c.cursorX++
	c.cursorY += c.cursorX / c.width
	c.cursorX = c.cursorX % c.width
	if c.cursorY >= c.top+c.height {
		c.scrollDown()
	}
}

// decrementCursor will move the cursor one cell to the left, scrolling the
// screen if necessary, but without redrawing the cursor.
func (c *Console) decrementCursor() {
	c.cursorX--
	if c.cursorX < 0 {
		c.cursorY--
		c.cursorX = c.left + c.width
	}
}

// moveCursorLeft will move the cursor one cell to the left (if possible) and
// will redraw the cursor image.
func (c *Console) moveCursorLeft() {
	if loc := c.getCursorLoc(); loc == 0 {
		return
	}
	tb.SetCell(c.cursorX, c.cursorY, c.getCursorChar(), tb.ColorDefault, tb.ColorDefault)
	c.decrementCursor()
	tb.SetCell(c.cursorX, c.cursorY, c.getCursorChar(), cursorFmt, cursorFmt)
}

// moveCursorLeft will move the cursor one cell to the right (if possible) and
// will redraw the cursor image.
func (c *Console) moveCursorRight() {
	if loc := c.getCursorLoc(); loc >= len(c.currline) {
		return
	}
	tb.SetCell(c.cursorX, c.cursorY, c.getCursorChar(), tb.ColorDefault, tb.ColorDefault)
	c.incrementCursor()
	tb.SetCell(c.cursorX, c.cursorY, c.getCursorChar(), cursorFmt, cursorFmt)
}

// getCursorLoc will return the offset into the line that the cursor is at.
// The prompt is not included, meaning that the first cell AFTER the prompt
// will have an offset of 0.
func (c *Console) getCursorLoc() int {
	return (c.cursorY-c.promptY)*c.width + c.cursorX - c.left - len(c.prompt)
}

// getCursorChar returns the character currently underneath the cursor.
func (c *Console) getCursorChar() rune {
	if loc := c.getCursorLoc(); loc < 0 || loc >= len(c.currline) {
		return ' '
	} else {
		return rune(c.currline[loc])
	}
}

// writeChar will write a character where the cursor is. It will NOT shift
// any characters that occur after it on the line it is on to the right.
func (c *Console) writeChar(ch rune) {
	tb.SetCell(c.cursorX, c.cursorY, ch, tb.ColorDefault, tb.ColorDefault)
	c.incrementCursor()
	tb.SetCell(c.cursorX, c.cursorY, ' ', tb.ColorWhite, tb.ColorWhite)
}

// insertChar will write a character where the cursor is, shifting all
// characters that occur after it on the line it is on to the right.
func (c *Console) insertChar(ch rune) {
	loc := c.getCursorLoc()
	cX := c.cursorX
	cY := c.cursorY
	// Remove cursor box from image.
	tb.SetCell(cX, cY, ch, tb.ColorDefault, tb.ColorDefault)
	// Move cursor to the right.
	c.incrementCursor()

	// Shift all later cells to the right.
	for i := loc; i < len(c.currline); i++ {
		tb.SetCell(c.cursorX, c.cursorY, rune(c.currline[i]), tb.ColorDefault, tb.ColorDefault)
		c.incrementCursor()
	}

	// Move the cursor back to where it was.
	c.cursorX = cX
	c.cursorY = cY
	// Update the current line.
	c.currline = c.currline[:loc] + string(ch) + c.currline[loc:]
	// Shift the cursor one over to be over the character AFTER the one
	// that was just inserted.
	c.moveCursorRight()
}

// Print prints a string, with no newline, to a given Console.
// This will have strange behavior unless called when the cursor is at the
// end of the current line.
func (c *Console) Print(str string) {
	for _, ch := range str {
		c.writeChar(ch)
	}
}

// Println prints a string, followed by a newline, to a given Console.
// This will have strange behavior unless called when the cursor is at the
// end of the current line.
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

// backspace performs the equivalent of pressing the backspace key.
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

// clearLine clears all text to the right of the current line's prompt string.
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
