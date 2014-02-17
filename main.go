package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	version = "0.1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s server [channels]\n", os.Args[0])

		os.Exit(1)
	}

	log.Printf("Launching raven (version: %s)\n", version)

	raven, err := NewRaven()

	if err != nil {
		log.Fatalf("Failed to create IRC client: %s\n", err)
	}

	defer raven.Disconnect()

	raven.LoadModules([]string{"core"})

	if len(os.Args) >= 3 {
		raven.Subscribe("welcome", func(raven *Raven, message *Message) {
			for _, channel := range strings.Split(os.Args[2], ",") {
				raven.Join(strings.TrimSpace(channel))
			}
		})
	}

	if err := raven.Connect(os.Args[1]); err == nil {
		log.Printf("Connected to %s\n", os.Args[1])
	} else {
		log.Fatalf("Failed to connect to %s\n", os.Args[1])
	}

	raven.Fly()

	log.Printf("Shutting down")
}
