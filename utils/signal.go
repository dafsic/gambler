package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignals = []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}

func InterruptListener() <-chan struct{} {
	c := make(chan struct{})
	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, interruptSignals...)

		// Listen for initial shutdown signal and close the returned
		// channel to notify the caller.
		select {
		case sig := <-interruptChannel:
			fmt.Printf("Received signal (%s).  Shutting down...\n", sig)
		}
		close(c)

		// Listen for repeated signals and display a message so the user
		// knows the shutdown is in progress and the process is not
		// hung.
		for {
			select {
			case sig := <-interruptChannel:
				fmt.Printf("Received signal (%s).  Already "+
					"shutting down...\n", sig)
			}
		}
	}()

	return c
}
