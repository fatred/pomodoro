package timer

import (
	"fmt"
	"strings"
	"time"
)

type DisplayCallback func(msg string, update bool)

type Pomodoro struct {
	WorkSec          int
	ShortBreakSec    int
	LongBreakSec     int
	CyclesBeforeLong int
	currentCycle     int
	sessionType      string
	callback         DisplayCallback
	lastMessage      string
	quit             chan struct{}
}

func New(work, short, long, cycles int) *Pomodoro {
	return &Pomodoro{
		WorkSec:          work * 60,
		ShortBreakSec:    short * 60,
		LongBreakSec:     long * 60,
		CyclesBeforeLong: cycles,
		sessionType:      "work",
		quit:             make(chan struct{}),
	}
}

func (p *Pomodoro) Start(display func(string, bool)) {
	if display == nil {
		return
	}
	p.callback = display
	p.sessionType = "work"

	go func() {
		for {
			select {
			case <-p.quit:
				return
			default:
				p.runSession()
				if !p.promptContinue() {
					p.callback("Goodbye!", false)
					return
				}
			}
		}
	}()
}

func (p *Pomodoro) runSession() {
	if p.callback == nil {
		return
	}
	var duration int
	switch p.sessionType {
	case "work":
		duration = p.WorkSec
	case "short":
		duration = p.ShortBreakSec
	default:
		duration = p.LongBreakSec
	}

	p.callback("--- "+p.sessionType+" START ---", false)
	for remaining := duration; remaining >= 0; remaining-- {
		mins, secs := remaining/60, remaining%60
		message := fmt.Sprintf("%02d:%02d", mins, secs)
		p.lastMessage = message // Store last message
		p.callback(message, true)
		time.Sleep(time.Second)
	}
	p.transition()
}

func (p *Pomodoro) transition() {
	if p.sessionType == "work" {
		p.currentCycle++
		if p.currentCycle%p.CyclesBeforeLong == 0 {
			p.sessionType = "long"
		} else {
			p.sessionType = "short"
		}
	} else {
		p.sessionType = "work"
	}
}

func (p *Pomodoro) promptContinue() bool {
	p.callback("Press 'c' to continue or any other key to quit.", false)
	var input string
	fmt.Scanln(&input)
	return strings.ToLower(input) == "c"
}

func (p *Pomodoro) RefreshDisplay() {
	if p.callback != nil && p.lastMessage != "" {
		p.callback(p.lastMessage, true)
	}
}

func (p *Pomodoro) Stop() {
	close(p.quit)
}

func (p *Pomodoro) LastMessage() string {
	return p.lastMessage
}
