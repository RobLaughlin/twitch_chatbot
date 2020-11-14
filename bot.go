package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/roblaughlin/twitch-chatbot/env"
	"github.com/roblaughlin/twitch-chatbot/irc"
)

const envpath string = "bot.env"
const commandPrefix = "~"

func messageHandler(client *irc.Client, message *irc.Message) error {
	switch message.Command {
	case "PING":
		return client.Pong()
	case "PRIVMSG":
		if strings.HasPrefix(message.Contents, commandPrefix) {
			args := strings.Split(strings.Split(message.Contents, commandPrefix)[1], " ")
			command := args[0]

			switch strings.ToLower(command) {
			case "ping":
				return client.Ping()
			case "sendsomething":
				if len(args) > 1 {
					return client.Send(args[1], message.Channel)
				}

			}
		}
	default:
	}
	fmt.Println(*message)
	return nil
}

func main() {
	// Make sure all client variables exist in client.env
	env, err := env.Validate(envpath, []string{
		"HOST",
		"PORT",
		"USER",
		"PASS",
	})

	// Check for valid port number
	port, err := strconv.Atoi(env["PORT"])
	if err != nil {
		log.Fatalf("Invalid port number: %s.", env["PORT"])
	}

	// Initialize a connection
	connection := irc.NewConnection(env["HOST"], (uint16)(port))
	bot := irc.NewClient(env["USER"], env["PASS"], connection)

	// Dial into the connection
	fmt.Printf("Dialing into %s...\n", connection.String())
	if !bot.Connected() {
		bot.Connect(true)
		bot.Join("#gravitybotv2")
	}

	for {
		_, err := bot.ParseMessage(messageHandler)
		if err != io.EOF {
			fmt.Println(err.Error())
		}
	}
}
