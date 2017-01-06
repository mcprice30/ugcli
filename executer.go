// ugCLI is a library built atop termbox for creating CLI applications.
package ugcli

// The Executer interface will in the future to allow
// users to power the commands behind an ugCLI console.

type Executer interface {
	Execute(command string) int
}
