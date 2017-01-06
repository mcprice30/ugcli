// ugCLI is a library built atop termbox for creating CLI applications.
package ugcli

// The Completer interface will be used in the future
// to allow tab-completion in an ugCLI console.

type Completer interface {
	Complete(input string) []string
}
