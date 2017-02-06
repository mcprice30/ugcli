// Package console defines an ugcli component that contains prebuilt
// functionality for building interactive prompts. A console consists of a
// prompt, where users can enter commands. In addition to simply allowing for
// commands to be typed, additional features such as pressing the arrow
// keys to view past commands are implemented.
//
// Additionally, additional functionality can be added to a console via the
// use of an executer, which will take and exectue a command run in the console,
// or a completer, which will implement tab completion within the console.
package console

// defaultPrompt indicates the default prefix to be displayed before all
// commands wihtin the console.
const defaultPrompt = "> "

// bufferSize indicates the maximum number of previously executed commands
// to buffer in terms of using the up/down arrows to view past commands.
const bufferSize = 100

// Console represents a console component of a command line application.
// Since it implements the component interface, it can be embedded into ugcli
// applications to deal with custom interations with the user.
type Console struct {

	// Which cell row of the terminal the console starts at.
	top int

	// Which cell column of the terminal the console starts at.
	left int

	// How many cell columns wide the console is.
	width int

	// How many cell rows tall the console is.
	height int

	// The current cell column the cursor is located at.
	// Note that this is indexed from 0 starting with the leftmost column of the
	// terminal window, NOT the top row of the console.
	cursorX int

	// The current cell row the cursor is located at.
	// Note that this is indexed from 0 starting with the top row of the terminal
	// window, NOT the top row of the console.
	cursorY int

	// What text is printed as the prompt.
	prompt string

	// What cell row the prompt was most recently located at.
	// Note that this is indexed from 0 starting with the top row of the terminal
	// window, NOT the top row of the console.
	promptY int

	// The text of the current line of the console (what the user is actively
	// editing).
	currline string

	// A user defined executer, used to process the actual commands sent to
	// the console.
	executer Executer

	// A user defined completer, used to provide suggestions for tab completion.
	completer Completer

	// How many lines up into the previous commands buffer the user currently is.
	// This is used when pressing the arrows to cycle through old commands.
	diff int

	// Buffers what the current line was before the user started pressing the
	// up/down arrow keys to examine previous commands, so that what they were
	// typing can be restored if they come all the way back.
	oldLineCopy string

	// Holds up to bufferSize previous commands.
	lineBuffer []string

	// What index of the buffer the current line would be written into.
	bufferIdx int

	// Indicates whether the console is actively running right now.
	running bool
}

// New console will take the location and size of a console (in cells) and
// return an appropriate console component.
//
// Note that top and left are 0-indexed.
func NewConsole(top, left, width, height int) *Console {

	return &Console{
		top:         top,
		left:        left,
		width:       width,
		height:      height,
		cursorX:     left,
		cursorY:     top,
		promptY:     top,
		prompt:      defaultPrompt,
		currline:    "",
		diff:        0,
		oldLineCopy: "",
		lineBuffer:  make([]string, bufferSize),
		bufferIdx:   0,
		running:     true,
	}
}

// SetExecuter attaches a user-defined command execution object to the console.
func (c *Console) SetExecuter(e Executer) {
	c.executer = e
}

// SetCompleter attatches a user-defined tab completion object to the console.
func (c *Console) SetCompleter(comp Completer) {
	c.completer = comp
}
