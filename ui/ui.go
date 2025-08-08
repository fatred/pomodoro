package ui

import (
	"github.com/fatred/pomodoro/timer"
	"github.com/nsf/termbox-go"
)

func Run(p *timer.Pomodoro, quit chan struct{}) error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	w, h := termbox.Size()

	display := func(msg string, update bool) {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		if update {
			drawCentered(msg, w/2, h/2)
		} else {
			drawString(msg, 1, 0)
		}
		termbox.Flush()
	}

	// Start the display loop
	p.Start(display)
	defer p.Stop()

	// Main event loop
	for {
		select {
		case <-quit:
			return nil
		default:
			ev := termbox.PollEvent()
			switch ev.Type {
			case termbox.EventKey:
				if ev.Ch == 'q' || ev.Ch == 'Q' || ev.Key == termbox.KeyCtrlC {
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					termbox.Flush()
					if confirmDialog() {
						return nil
					}
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					display(p.LastMessage(), true)
				}
			case termbox.EventError:
				return ev.Err
			}
		}
	}
}

func drawCentered(msg string, x, y int) {
	x -= len(msg) / 2
	drawString(msg, x, y)
}

func drawString(msg string, x, y int) {
	for i, c := range msg {
		termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}
