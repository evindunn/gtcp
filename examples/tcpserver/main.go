package main

import (
	"fmt"
	logger "github.com/evindunn/gologyourself"
	"github.com/evindunn/gtcp/pkg/tcpmessage"
	"github.com/evindunn/gtcp/pkg/tcpserver"
	"net"
	"os"
	"strconv"
)

// Handler is a type implementing the github.com/evindunn/gtcp/pkg/tcpserver/ConnectionHandler interface
type Handler struct{
	server *tcpserver.Server
	srvLogger logger.Logger
}

// HandleConnection handles a connection for github.com/evindunn/gtcp/pkg/tcpserver
func (h *Handler) HandleConnection(c *net.Conn) {
	connection := *c
	defer connection.Close()

	client := connection.RemoteAddr().String()
	msg, err := tcpmessage.MessageFromConnection(c)

	if err != nil {
		h.srvLogger.Log(logger.LevelError, fmt.Sprintf("[%s] Error parsing connection: %s", client, err))
	} else {
		h.srvLogger.Log(
			logger.LevelInfo,
			fmt.Sprintf(
				"[%s] Size: %d, Compressed: %v, Content: %s", client,
				msg.GetSize(),
				msg.IsCompressed(),
				string(msg.GetContent())),
		)

		sendMsg := tcpmessage.NewMessage("PONG", false)
		_, err := connection.Write(sendMsg.ToBytes())
		if err != nil {
			h.srvLogger.Log(
				logger.LevelError,
				fmt.Sprintf("[%s] Error writing to connection: %s", client, err),
			)
		}
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
	h.srvLogger = logger.NewLogger(logger.LevelDebug)
	srv, err := tcpserver.NewServer(port, &h)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred creating the server: %s\n", err)
		os.Exit(1)
	}

	srv.Start()
}
