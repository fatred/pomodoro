package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatred/pomodoro/timer"
	"github.com/fatred/pomodoro/ui"
)

func main() {
	work := flag.Int("w", 25, "Work duration in minutes")
	short := flag.Int("s", 5, "Short break duration in minutes")
	long := flag.Int("l", 15, "Long break duration in minutes")
	cycles := flag.Int("c", 4, "Work sessions before long break")
	flag.Parse()

	p := timer.New(*work, *short, *long, *cycles)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create quit channel
	quitChan := make(chan struct{})

	// Handle interrupt in a separate goroutine
	go func() {
		<-sigChan
		log.Println("\nReceived interrupt signal. Exiting...")
		close(quitChan)
	}()

	if err := ui.Run(p, quitChan); err != nil {
		log.Fatal(err)
	}
}
