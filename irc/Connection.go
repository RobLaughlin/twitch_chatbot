package irc

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
)

// Connection stream to a server
type Connection struct {
	host      string
	port      uint16
	connected bool
	stream    net.Conn
	reader    *textproto.Reader
}

// NewConnection create a new connection to a server
func NewConnection(host string, port uint16) *Connection {
	return &Connection{
		host:      host,
		port:      port,
		connected: false,
		stream:    nil,
		reader:    nil,
	}
}

// Connect to a given connection stream
func (connection *Connection) Connect() (net.Conn, error) {
	if connection.connected {
		connection.stream.Close()
		connection.reader = nil
	}

	conn, err := net.Dial("tcp", connection.String())

	if err != nil {
		connection.connected = false
		connection.stream = nil
		connection.reader = nil
		return nil, err
	}

	connection.stream = conn
	connection.reader = textproto.NewReader(bufio.NewReader(connection.stream))
	connection.connected = true

	return conn, nil
}

// Send data to the server connection stream
func (connection Connection) Send(message string) error {
	if connection.connected {
		_, err := fmt.Fprintf(connection.stream, message)
		return err
	}

	return fmt.Errorf("cannot send message {%s} on disconnected stream {%s}", message, connection.String())
}

// ReadLine from the connection stream
func (connection Connection) ReadLine() (string, error) {
	return connection.reader.ReadLine()
}

// Connected - Get the connection status of the stream
func (connection Connection) Connected() bool {
	return connection.connected
}

// String together the host and the port
func (connection Connection) String() string {
	return fmt.Sprintf("%s:%d", connection.host, connection.port)
}
