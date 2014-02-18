package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

func init() {
	RegisterModule("core", func() Module {
		return &Core{}
	})
}

type Core struct{}

func (core *Core) Init(client *IRCClient) error {
	client.Subscribe("connected", core.auth)
	client.Subscribe("ping", core.pong)
	client.Subscribe("nicknameinuse", core.changeNick)

	client.Subscribe("!quit", core.quit)
	client.Subscribe("!join", core.joinChannel)
	client.Subscribe("!part", core.partChannel)

	return nil
}

func (core *Core) auth(client *IRCClient, message *Message) {
	fmt.Fprintf(client.connection, "USER client 0 0 :client\r\n")
	fmt.Fprintf(client.connection, "NICK client\r\n")
}

func (core *Core) pong(client *IRCClient, message *Message) {
	fmt.Fprintf(client.connection, "PONG :%s\r\n", message.args[0])
}

func (core *Core) changeNick(client *IRCClient, message *Message) {
	nick := fmt.Sprintf("client-%d", rand.Intn(999))

	log.Printf("Changing nick to: %s\n", nick)
	fmt.Fprintf(client.connection, "NICK %s\r\n", nick)
}

func (core *Core) quit(client *IRCClient, message *Message) {
	/* TODO: code me */
}

func (core *Core) joinChannel(client *IRCClient, message *Message) {
	for _, channel := range strings.Split(message.args[1], ",") {
		client.Join(strings.TrimSpace(channel))
	}
}

func (core *Core) partChannel(client *IRCClient, message *Message) {
	for _, channel := range strings.Split(message.args[1], ",") {
		client.Part(strings.TrimSpace(channel))
	}
}
