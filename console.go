// ugCLI is a library built atop termbox for creating CLI applications.
package ugcli

// Consoles are an ugCLI object that contains prebuilt
// functionality for building interactive prompts.
// In the future, custom functionality will be added
// using Completers and Executers

const default_prompt = "> "

// The Console struct
type Console struct {
	top     int
	left    int
	width   int
	height  int
	cursorX int
	cursorY int
	prompt  string

	executer  Executer
	completer Completer
}

// Create a new console, with the given coordinates.
func NewConsole(top, left, width, height int) *Console {

	return &Console{
		top:     top,
		left:    left,
		width:   width,
		height:  height,
		cursorX: left,
		cursorY: top,
		prompt:  default_prompt,
	}
}

func (c *Console) SetExecuter(e Executer) {
	c.executer = e
}

func (c *Console) SetCompleter(comp Completer) {
	c.completer = comp
}
