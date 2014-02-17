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

func (core *Core) Init(raven *Raven) error {
	raven.Subscribe("connected", core.auth)
	raven.Subscribe("ping", core.pong)
	raven.Subscribe("nicknameinuse", core.changeNick)

	raven.Subscribe("!quit", core.quit)
	raven.Subscribe("!join", core.joinChannel)
	raven.Subscribe("!part", core.partChannel)

	return nil
}

func (core *Core) auth(raven *Raven, message *Message) {
	fmt.Fprintf(raven.connection, "USER raven 0 0 :raven\r\n")
	fmt.Fprintf(raven.connection, "NICK raven\r\n")
}

func (core *Core) pong(raven *Raven, message *Message) {
	fmt.Fprintf(raven.connection, "PONG :%s\r\n", message.args[0])
}

func (core *Core) changeNick(raven *Raven, message *Message) {
	nick := fmt.Sprintf("raven-%d", rand.Intn(999))

	log.Printf("Changing nick to: %s\n", nick)
	fmt.Fprintf(raven.connection, "NICK %s\r\n", nick)
}

func (core *Core) quit(raven *Raven, message *Message) {
	/* TODO: code me */
}

func (core *Core) joinChannel(raven *Raven, message *Message) {
	for _, channel := range strings.Split(message.args[1], ",") {
		raven.Join(strings.TrimSpace(channel))
	}
}

func (core *Core) partChannel(raven *Raven, message *Message) {
	for _, channel := range strings.Split(message.args[1], ",") {
		raven.Part(strings.TrimSpace(channel))
	}
}
