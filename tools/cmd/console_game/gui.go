package main

import (
	"fmt"
	"log"

	c "github.com/jroimartin/gocui"
	"github.com/pkg/errors"
)

const (
	// List box width.
	lw = 20
	// Input box height.
	ih = 3
	// Debug lines height
	db = 4
)

// Items to fill the list with.
var listItems = []string{
	"Health: 20",
	"Cards inHand",
	"---------",
	"Minion",
	"General",
	"Bob",
}

func runGocui() {
	// Create a new GUI.
	g, err := c.NewGui(c.OutputNormal)
	if err != nil {
		log.Println("Failed to create a GUI:", err)
		return
	}
	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(layout)

	err = g.SetKeybinding("", c.KeyCtrlC, c.ModNone, quit)
	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}

	// Now let's define the views.

	// The terminal's width and height are needed for layout calculations.
	tw, th := g.Size()

	// First, create the list view.
	lv, err := g.SetView("player1", 0, 0, lw, (th/2)-1)
	// ErrUnknownView is not a real error condition.
	// It just says that the view did not exist before and needs initialization.
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create main view:", err)
		return
	}
	lv.Title = "Player1"
	lv.FgColor = c.ColorCyan

	// First, create the player2 list view.
	lv2, err := g.SetView("player2", 0, (th / 2), lw, th-1)
	// ErrUnknownView is not a real error condition.
	// It just says that the view did not exist before and needs initialization.
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create main view:", err)
		return
	}
	lv2.Title = "Player2"
	lv2.FgColor = c.ColorGreen

	//---------------board
	// Then the output view.
	ov, err := g.SetView("board", lw+1, 0, tw-1, th-db-ih-1)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create board view:", err)
		return
	}
	ov.Title = "ZugB - Console edition - Play Board"
	ov.FgColor = c.ColorGreen
	// Let the view scroll if the output exceeds the visible area.
	ov.Autoscroll = true

	// Then the output view.
	bp1, err := g.SetView("board-player1", lw+2, 1, tw-3, th-db-ih-4)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create board view:", err)
		return
	}
	bp1.Title = "Player 1 - Board"
	bp1.FgColor = c.ColorGreen
	// Let the view scroll if the output exceeds the visible area.
	bp1.Autoscroll = true
	_, err = fmt.Fprintln(bp1, "Cards on Table\n General #1 \n General #2, Minion")
	check(err)

	player1Height := th - db - ih - 3
	bp2, err := g.SetView("board-player2", lw+2, player1Height/2, tw-3, th-db-ih-2)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create board view:", err)
		return
	}
	bp2.Title = "Player 2 - Board"
	bp2.FgColor = c.ColorRed
	// Let the view scroll if the output exceeds the visible area.
	bp2.Autoscroll = true
	_, err = fmt.Fprintln(bp2, "Cards on Table\n General #1 \n General #2, Minion")
	check(err)

	///------------end board

	db, err := g.SetView("Debug Logs", lw+1, th-db-ih-1, tw-1, th-1)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create board view:", err)
		return
	}
	db.Title = "Debug Logs"
	db.FgColor = c.ColorGreen
	// Let the view scroll if the output exceeds the visible area.
	db.Autoscroll = true
	_, err = fmt.Fprintln(db, "Incoming Transactions ....")
	if err != nil {
		log.Println("Failed to print into output view:", err)
	}
	_, err = fmt.Fprintln(db, "Press Ctrl-c to quit")
	if err != nil {
		log.Println("Failed to print into output view:", err)
	}

	// And finally the input view.
	iv, err := g.SetView("input", lw+1, th-ih, tw-1, th-1)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create input view:", err)
		return
	}
	iv.Title = "Input"
	iv.FgColor = c.ColorYellow
	// The input view shall be editable.
	iv.Editable = true
	err = iv.SetCursor(0, 0)
	if err != nil {
		log.Println("Failed to set cursor:", err)
		return
	}

	// Make the enter key copy the input to the output.
	err = g.SetKeybinding("input", c.KeyEnter, c.ModNone, func(g *c.Gui, iv *c.View) error {
		// We want to read the view's buffer from the beginning.
		iv.Rewind()

		// Get the output view via its name.
		ov, e := g.View("board")
		if e != nil {
			log.Println("Cannot get board view:", e)
			return e
		}
		// Thanks to views being an io.Writer, we can simply Fprint to a view.
		_, e = fmt.Fprint(ov, iv.Buffer())
		if e != nil {
			log.Println("Cannot print to output view:", e)
		}
		// Clear the input view
		iv.Clear()
		// Put the cursor back to the start.
		e = iv.SetCursor(0, 0)
		if e != nil {
			log.Println("Failed to set cursor:", e)
		}
		return e

	})
	if err != nil {
		log.Println("Cannot bind the enter key:", err)
	}

	// Fill the list view.
	for _, s := range listItems {
		// Again, we can simply Fprint to a view.
		_, err = fmt.Fprintln(lv, s)
		if err != nil {
			log.Println("Error writing to the list view:", err)
			return
		}
		_, err = fmt.Fprintln(lv2, s)
		if err != nil {
			log.Println("Error writing to the list view:", err)
			return
		}
	}

	// Set the focus to the input view.
	_, err = g.SetCurrentView("input")
	if err != nil {
		log.Println("Cannot set focus to input view:", err)
	}

	// Start the main loop.
	err = g.MainLoop()
	log.Println("Main loop has finished:", err)
}

// The layout handler calculates all sizes depending
// on the current terminal size.
func layout(g *c.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()

	// Update the views according to the new terminal size.
	_, err := g.SetView("player1", 0, 0, lw, (th/2)-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update player1 view")
	}
	_, err = g.SetView("player2", 0, (th / 2), lw, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update player1 view")
	}
	_, err = g.SetView("board", lw+1, 0, tw-1, th-ih-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update board view")
	}
	_, err = g.SetView("Debug Logs", lw+1, th-db-ih-1, tw-1, th-1)
	if err != nil && err != c.ErrUnknownView {
		log.Println("Failed to create board view:", err)
		return errors.Wrap(err, "Cannot update debug view")
	}
	_, err = g.SetView("board-player1", lw+2, 1, tw-3, th-db-ih-4)
	if err != nil {
		return errors.Wrap(err, "Cannot update board view")
	}

	player1Height := th - db - ih - 3
	_, err = g.SetView("board-player2", lw+2, player1Height/2, tw-3, th-db-ih-2)
	if err != nil {
		return errors.Wrap(err, "Cannot update board view")
	}
	_, err = g.SetView("input", lw+1, th-ih, tw-1, th-1)
	if err != nil {
		return errors.Wrap(err, "Cannot update input view.")
	}
	return nil
}

// `quit` is a handler that gets bound to Ctrl-C.
// It signals the main loop to exit.
func quit(g *c.Gui, v *c.View) error {
	return c.ErrQuit
}

/*
Our main func just needs to read the name from the TUI lib from the command line
and execute the respective code.
*/

// would never use this in a serious program
func check(err error) {
	if err != nil {
		panic(err)
	}
}
