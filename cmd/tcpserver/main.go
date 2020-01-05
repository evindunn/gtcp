package main

import (
	"fmt"
	"github.com/evindunn/gtcp/pkg/tcpmessage"
	"github.com/evindunn/gtcp/pkg/tcpserver"
	"log"
	"net"
	"os"
	"strconv"
)

type Handler struct{}

func (h *Handler) HandleConnection(c *net.Conn) {
	connection := *c
	msg, err := tcpmessage.MessageFromConnection(c)
	if err != nil {
		log.Printf("[%s] Error parsing connection: %s\n", (*c).RemoteAddr().String(), err)
	} else {
		log.Printf(
			"[%s] Size: %d, Compressed: %v, Content: %s", connection.RemoteAddr().String(),
			msg.GetSize(),
			msg.IsCompressed(),
			string(msg.GetContent()))
	}

	err = connection.Close()
	if err != nil {
		log.Printf("[%s] Error closing connection connection: %s\n", (*c).RemoteAddr().String(), err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Specify a port")
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Port must be an integer")
		os.Exit(1)
	}

	var h Handler
	srv, err := tcpserver.NewServer(port, &h)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred creating the server: %s\n", err)
		os.Exit(1)
	}

	srv.Start()
}
