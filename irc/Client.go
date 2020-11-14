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
	return err
}

// Join & authenticate using oauth with the given connection
func (client Client) Join(oauth bool) error {
	if !client.Connection.connected {
		if err := client.Connect(); err != nil {
			return err
		}
	}

	var userAuth string = fmt.Sprintf("NICK %s\n", client.username)
	var passAuth string

	if oauth {
		passAuth = fmt.Sprintf("PASS oauth:%s\n", client.password)
	} else {
		passAuth = fmt.Sprintf("PASS %s\n", client.password)
	}

	if err := client.Connection.Send(passAuth); err != nil {
		return err
	}
	if err := client.Connection.Send(userAuth); err != nil {
		return err
	}

	return nil
}
