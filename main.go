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

	client, err := NewIRCClient()

	if err != nil {
		log.Fatalf("Failed to create IRC client: %s\n", err)
	}

	defer client.Disconnect()

	client.LoadModules([]string{"core"})

	if len(os.Args) >= 3 {
		client.Subscribe("welcome", func(client *IRCClient, message *Message) {
			for _, channel := range strings.Split(os.Args[2], ",") {
				client.Join(strings.TrimSpace(channel))
			}
		})
	}

	if err := client.Connect(os.Args[1]); err == nil {
		log.Printf("Connected to %s\n", os.Args[1])
	} else {
		log.Fatalf("Failed to connect to %s\n", os.Args[1])
	}

	client.Run()

	log.Printf("Shutting down")
}
