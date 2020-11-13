package irc

import "fmt"

type Client struct {
	username   string ""
	password   string ""
	Connection *Connection
}

func NewClient(username string, password string, connection *Connection) *Client {
	return &Client{
		username:   username,
		password:   password,
		Connection: connection,
	}
}

func (client Client) Connect() error {
	_, err := client.Connection.Connect()
	return err
}

func (client Client) Join() error {
	if !client.Connection.connected {
		if err := client.Connect(); err != nil {
			return err
		}
	}

	if err := client.Connection.Send(fmt.Sprintf("PASS oauth:%s\n", client.password)); err != nil {
		return err
	}
	if err := client.Connection.Send(fmt.Sprintf("NICK %s\n", client.username)); err != nil {
		return err
	}

	return nil
}
