package ui

import "github.com/nsf/termbox-go"

func confirmDialog() bool {
	w, h := termbox.Size()
	msg := "Are you sure you want to quit? (y/n)"

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	drawCentered(msg, w/2, h/2)
	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 'y' || ev.Ch == 'Y' {
				return true
			}

			if ev.Ch == 'n' || ev.Ch == 'N' {
				return false
			}
		}
	}
}
