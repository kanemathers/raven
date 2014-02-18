package main

import "fmt"

func init() {
	RegisterModule("auth", func() Module {
		return &Auth{}
	})
}

type Auth struct{}

func (auth *Auth) Init(client *IRCClient) error {
	client.Subscribe("privmsg", auth.RequireAuth(auth.test))

	return nil
}

func (auth *Auth) IsAuthed(message *Message) bool {
	return false
}

func (auth *Auth) test(client *IRCClient, message *Message) {
	fmt.Println("here")
}

func (auth *Auth) RequireAuth(fn Handler) Handler {
	return func(client *IRCClient, message *Message) {
		if auth.IsAuthed(message) {
			fn(client, message)
		}
	}
}
