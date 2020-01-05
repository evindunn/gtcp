package main

import (
	"fmt"
	"github.com/evindunn/gtcp/pkg/tcpMessage"
	"github.com/evindunn/gtcp/pkg/tcpServer"
	"log"
	"net"
	"os"
	"strconv"
)

type Handler struct {}

func (h *Handler) HandleConnection(c *net.Conn) {
	connection := *c
	msg, err := tcpMessage.MessageFromConnection(c)
	if err != nil {
		log.Printf("[%s] Error parsing connection: %s\n", (*c).RemoteAddr().String(), err)
	} else {
		log.Printf("[%s] %s", connection.RemoteAddr().String(), string(msg.GetContent()))
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
	srv, err := tcpServer.NewServer(port, &h)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred creating the server: %s\n", err)
		os.Exit(1)
	}

	srv.Start()
}