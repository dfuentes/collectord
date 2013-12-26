package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

func main() {
	flag.Parse()

	SetupConfig()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, os.Kill)

loop:
	for {
		select {
		case s := <-sigChannel:
			log.Printf("Received signal: %v", s)
			break loop
		}
	}
}
