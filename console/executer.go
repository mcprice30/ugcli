package console

// Executer allows ugcli users to specify their own rules for executing
// commands within the console.
type Executer interface {
	// Execute will take a command passed from the console, perform some action,
	// which may potentially include writing to the console, and will then
	// return an exit code for the command, along with indicating whether the
	// console should continue executing.
	Execute(command string) (statusCode int, keepRunning bool)

	// BoundConsole returns whatever console this executer will execute commands
	// for.
	BoundConsole() *Console
}

// DefaultExecuter will produce a builtin executer that supports the "exit"
// command, and otherwise prints out whatever was just executed. The executer
// will be bound to whatever console it was created with.
func DefaultExecuter(c *Console) Executer {
	return &echoExecuter{
		con: c,
	}
}

// echoExecuter is a default executer implementation that simply echoes text.
type echoExecuter struct {
	// con stores whatever console this executer is bound to.
	con *Console
}

// Execute will take a command and echo back that command, along with a
// 0 status code, unless the given command is exit, at which point it will
// inform the console this is bound to to stop execution.
func (ex *echoExecuter) Execute(command string) (statusCode int, keepRunning bool) {
	if command == "exit" {
		return 0, false
	} else if len(command) > 0 {
		ex.con.Println(command)
	}
	return 0, true
}

// BoundConsole returns whatever executer this executer is bound to.
func (ex *echoExecuter) BoundConsole() *Console {
	return ex.con
}
