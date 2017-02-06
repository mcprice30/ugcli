// This is a simple demo of the project thus far.
// In this demo, we initiate a termbox session,
// and launch a console that takes up the initial size
// of the terminal.
package main


import (
	"os"

	tb "github.com/nsf/termbox-go"

	"github.com/mcprice30/ugcli"
	"github.com/mcprice30/ugcli/console"
)

// main is the entry point for the application.
func main() {

	// Attempt to initialize the termbox.
	// TODO: contain inside ugcli run method.
	if err := tb.Init(); err != nil {
		os.Exit(1)
	}

	// Clear the termbox.
	if err := tb.Clear(tb.ColorDefault, tb.ColorDefault); err != nil {
		panic(err)
	}

	// Create a new console taking the size of the terminal.
	w, h := tb.Size()
	con := console.NewConsole(0, 0, w, h)

	// prefix tree completer.
	completer := console.NewListCompleter([]string{"a", "ab", "abc", "bad",
		"carrot", "jane", "jack"})
	con.SetCompleter(completer)

	// Initialize ugcli application and add console.
	cli := ugcli.NewCli()
	cli.AddComponent(con)

	// Launch the application.
	cli.Run()

	// Close the termbox session when done.
	tb.Close()
}
