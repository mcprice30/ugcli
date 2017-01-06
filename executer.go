// ugCLI is a library built atop termbox for creating CLI applications.
package ugcli

// The Executer interface will in the future to allow
// users to power the commands behind an ugCLI console.

type Executer interface {
	Execute(command string) (statusCode int, keepRunning bool)
	BoundConsole() *Console
}

func DefaultExecuter(c *Console) Executer {
	return &echoExecuter{
		con: c,
	}
}

type echoExecuter struct {
	con *Console
}

func (ex *echoExecuter) Execute(command string) (statusCode int, keepRunning bool) {
	if command == "exit" {
		return 0, false
	} else if len(command) > 0 {
		ex.con.Println(command)
	}
	return 0, true
}

func (ex *echoExecuter) BoundConsole() *Console {
	return ex.con
}
