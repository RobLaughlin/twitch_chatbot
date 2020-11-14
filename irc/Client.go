package irc

import (
	"fmt"
	"strings"
)

// Client IRC client
type Client struct {
	username   string ""
	password   string ""
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
func (client Client) Connect() error {
	_, err := client.Connection.Connect()
	if err != nil {
		return err
	}
	return nil
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
func (client Client) Join(verbose bool) error {
	if !client.Connection.Connected() {
		if err := client.Connect(); err != nil {
			return err
		}
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
		fmt.Printf("User %s successfully joined: %s!\n", client.username, client.Connection.String())
	}

	return nil
}
