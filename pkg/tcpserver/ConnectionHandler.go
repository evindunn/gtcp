package tcpserver

import (
	"net"
)

// ConnectionHandler defines an interface for handling connections with a tcpServer.Server
type ConnectionHandler interface {
	HandleConnection(c *net.Conn)
}
