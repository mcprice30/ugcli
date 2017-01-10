// This directory contains demos of the project.
package main

// This is a simple demo of the project thus far.
// In this demo, we initiate a termbox session,
// and launch a console that takes up the initial size
// of the terminal.

import (
	"os"

	tb "github.com/nsf/termbox-go"

	"github.com/mcprice30/ugcli"
)

func main() {

	// Attempt to initialize the termbox.
	if err := tb.Init(); err != nil {
		os.Exit(1)
	}

	// Clear the termbox.
	if err := tb.Clear(tb.ColorDefault, tb.ColorDefault); err != nil {
		panic(err)
	}

	// Create a new console taking the size of the terminal.
	w, h := tb.Size()
	con := ugcli.NewConsole(0, 0, w, h)

	completer := ugcli.CreateListCompleter([]string{"a", "ab", "abc", "bad",
		"carrot", "jane", "jack"})

	con.SetCompleter(completer)

	// Launch the terminal.
	con.Run()

	// Close the termbox session when done.
	tb.Close()
}
