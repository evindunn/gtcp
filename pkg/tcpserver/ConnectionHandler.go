package tcpserver

import (
	"net"
)

type ConnectionHandler interface {
	HandleConnection(c *net.Conn)
}
