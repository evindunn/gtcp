package tcpServer

import (
	"net"
)

type ConnectionHandler interface {
	HandleConnection(c *net.Conn)
}