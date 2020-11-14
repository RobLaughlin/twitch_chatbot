package irc

import (
	"errors"
	"fmt"
	"strings"
)

// ParseHandler Handler to pass parsed data to after sanitization
type ParseHandler func(client *Client, message *Message) error

// Message parsed from IRC reader
type Message struct {
	Sender   string
	Command  string
	Channel  string
	Contents string
}

// NewMessage Message constructor
func NewMessage(sender string, status string, channel string, contents string) *Message {
	return &Message{
		Sender:   sender,
		Command:  status,
		Channel:  channel,
		Contents: contents,
	}
}

// Client IRC client
type Client struct {
	username   string
	password   string
	Connection *Connection
}

// NewClient Create a new client with a username, password (optional oauth token), and a connection.
func NewClient(username string, password string, connection *Connection) *Client {
	return &Client{
		username:   username,
		password:   password,
		Connection: connection,
	}
}

// Connect to the set connection
func (client Client) Connect(verbose bool) error {
	_, err := client.Connection.Connect()
	if err != nil {
		return err
	}
	return client.authenticate(verbose)
}

// Disconnect from the set connection
func (client Client) Disconnect() error {
	return client.Connection.Disconnect()
}

// Connected Check if client is connected
func (client Client) Connected() bool {
	return client.Connection.Connected()
}

// Join the connected IRC stream with the given nickname and password
func (client Client) authenticate(verbose bool) error {
	if !client.Connected() {
		return errors.New("Cannot authenticate to an unconnected stream")
	}

	var userAuth string = fmt.Sprintf("NICK %s\n", client.username)
	var passAuth string = fmt.Sprintf("PASS %s\n", client.password)
	var asteriskedPass string = strings.Repeat("*", len(client.password))

	if verbose {
		fmt.Printf("Joining %s with user: %s\n", client.Connection.String(), client.username)
		fmt.Printf("Sending PASS %s\n", asteriskedPass)
	}

	if err := client.Connection.Send(passAuth); err != nil {
		if verbose {
			fmt.Printf("There was an error sending PASS: %s\n", asteriskedPass)
		}
		return err
	}

	if verbose {
		fmt.Printf("Sending USER %s\n", client.username)
	}

	if err := client.Connection.Send(userAuth); err != nil {
		if verbose {
			fmt.Printf("There was an error sending USER: %s\n", client.username)
		}
		return err
	}

	if verbose {
		fmt.Printf("User %s successfully connected to: %s!\n", client.username, client.Connection.String())
	}

	return nil
}

// ParseMessage Parse an IRC message by reading a line from the connection reader.
// Returns a message object consisting of message information such as user, status, channel, and message.
// Returned only after callback finishes executing.
func (client Client) ParseMessage(callback ParseHandler) (*Message, error) {
	if client.Connected() {
		status, err := client.Connection.reader.ReadLine()
		if err != nil {
			return nil, err
		}
		fmt.Println(status)
		msg := NewMessage("", "", "", "")
		parts := strings.Split(status, ":")

		if len(parts) > 1 {

			// Check for the optional comamnd edge cases early.
			optCommand := strings.Trim(parts[0], " ")
			if optCommand != "" {
				msg.Command = optCommand
				msg.Contents = parts[1]
				callback(&client, msg)
				return msg, nil
			}

			// Split the prefix into parts separated by spaces.
			prefixParts := strings.Split(parts[1], " ")

			if len(prefixParts) > 2 {

				// Check for a valid username to parse. If there is no valid username, leave it empty.
				username := strings.Split(prefixParts[0], "!")
				if len(username) > 1 {
					msg.Sender = username[0]
				}

				// Assign a status: ex JOIN or PRIVMSG.
				msg.Command = prefixParts[1]

				// Assign a channel if it exists. Only assign a channel if the channel is not equal to the client username.
				if prefixParts[2] != strings.Trim(client.username, " ") {
					msg.Channel = strings.Trim(prefixParts[2], " ")
				}
			}

			// If there is an actual message to parse, assign it to contents.
			if len(parts) > 2 {
				msg.Contents = parts[2]
			}
		}

		return msg, callback(&client, msg)
	}
	return nil, errors.New("Client not connected, cannot parse message")

}

// Join a given IRC channel. Must include # prefix to indiciate it is in fact an IRC channel.
func (client Client) Join(channel string) error {
	return client.Connection.Send(fmt.Sprintf("JOIN %s", channel))
}

// Send a message to a given channel
func (client Client) Send(message string, channel string) error {
	return client.Connection.Send(fmt.Sprintf("PRIVMSG %s %s", channel, message))
}

// Pong Respond to a ping with a pong
func (client Client) Pong() error {
	return client.Connection.Send("PONG")
}

// Ping the connection stream
func (client Client) Ping() error {
	return client.Connection.Send("PING")
}
