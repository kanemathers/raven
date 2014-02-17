package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net"
	"strings"
	"time"
)

var (
	ircEvents = map[string]string{
		"001": "welcome",
		"433": "nicknameinuse",
	}
)

type Handler func(*Raven, *Message)

type Raven struct {
	connection net.Conn
	db         *sql.DB
	modules    map[string]Module
	handlers   map[string][]Handler
}

type Message struct {
	time    time.Time
	prefix  string
	command string
	args    []string
}

func NewRaven() (*Raven, error) {
	raven := &Raven{
		modules:  make(map[string]Module),
		handlers: make(map[string][]Handler),
	}

	if db, err := sql.Open("sqlite3", "/tmp/raven.db"); err == nil {
		raven.db = db
	} else {
		return nil, err
	}

	return raven, nil
}

func (raven *Raven) LoadModule(name string) error {
	if module := LoadModule(name); module != nil {
		log.Printf("Loading module: %s\n", name)

		raven.modules[name] = module

		if err := module.Init(raven); err != nil {
			log.Printf("Failed to load module: %s\n", name)
		}
	} else {
		log.Printf("Can't find module: %s\n", name)
	}

	return nil
}

func (raven *Raven) LoadModules(modules []string) error {
	for _, module := range modules {
		raven.LoadModule(module)
	}

	return nil
}

func (raven *Raven) Connect(server string) error {
	if connection, err := net.Dial("tcp", server); err == nil {
		raven.connection = connection
	} else {
		return err
	}

	raven.Dispatch("connected", nil)

	return nil
}

func (raven *Raven) Disconnect() {
	raven.connection.Close()
}

func (raven *Raven) Fly() error {
	reader := bufio.NewReader(raven.connection)

	for {
		line, err := reader.ReadString('\r')

		if err != nil {
			return err
		}

		prefix, command, args := parseLine(line)

		message := Message{
			time:    time.Now(),
			prefix:  prefix,
			command: command,
			args:    args,
		}

		raven.Dispatch("*", &message)

		if event, ok := ircEvents[command]; ok {
			raven.Dispatch(event, &message)
		} else {
			raven.Dispatch(strings.ToLower(command), &message)
		}

		if len(args) >= 2 && args[1][0] == '!' {
			s := strings.SplitN(args[1], " ", 2)

			if len(s) == 2 {
				message.args[1] = s[1]
			} else {
				message.args[1] = ""
			}

			raven.Dispatch(s[0], &message)
		}

		//fmt.Printf("%+v\n", message)
	}

	return nil
}

func (raven *Raven) Dispatch(event string, message *Message) {
	for _, handler := range raven.handlers[event] {
		handler(raven, message)
	}
}

func (raven *Raven) Subscribe(event string, fn Handler) {
	raven.handlers[event] = append(raven.handlers[event], fn)
}

func (raven *Raven) Join(channel string) {
	var buffer bytes.Buffer

	if channel[0] != '#' {
		buffer.WriteByte('#')
	}

	buffer.WriteString(channel)

	channel = buffer.String()

	log.Printf("Joining channel: %s\n", channel)
	fmt.Fprintf(raven.connection, "JOIN %s\r\n", channel)
}

func (raven *Raven) Part(channel string) {
	var buffer bytes.Buffer

	if channel[0] != '#' {
		buffer.WriteByte('#')
	}

	buffer.WriteString(channel)

	channel = buffer.String()

	log.Printf("Parting channel: %s\n", channel)
	fmt.Fprintf(raven.connection, "PART %s\r\n", channel)
}

func parseLine(line string) (prefix string, command string, args []string) {
	line = strings.Trim(line, "\r\n")

	if line[0] == ':' {
		s := strings.SplitN(line[1:], " ", 2)
		prefix = s[0]
		line = s[1]
	}

	if strings.Contains(line, " :") {
		s := strings.SplitN(line, " :", 2)
		args = strings.Split(s[0], " ")
		args = append(args, s[1])
	} else {
		args = strings.Split(line, " ")
	}

	command = args[0]
	args = append(args[:0], args[1:]...)

	return prefix, command, args
}
