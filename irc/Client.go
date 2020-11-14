package irc

import "fmt"

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

// Join & authenticate using oauth with the given connection
func (client Client) Join(oauth bool) error {
	if !client.Connection.Connected() {
		if err := client.Connect(); err != nil {
			return err
		}
	}

	var userAuth string = fmt.Sprintf("NICK %s\n", client.username)
	var passAuth string = fmt.Sprintf("PASS %s\n", client.password)

	if err := client.Connection.Send(passAuth); err != nil {
		return err
	}
	if err := client.Connection.Send(userAuth); err != nil {
		return err
	}

	return nil
}
