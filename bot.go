package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/roblaughlin/twitch-chatbot/env"
	"github.com/roblaughlin/twitch-chatbot/irc"
)

const envpath string = "bot.env"

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
	if !bot.Connection.Connected() {
		bot.Connect()
		bot.Join(true)
	}

	for {
		status, err := bot.Connection.ReadLine()
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}

		fmt.Println(status)
	}
}
