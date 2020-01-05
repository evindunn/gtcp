package main

import (
	"fmt"
	logger "github.com/evindunn/gologyourself"
	"github.com/evindunn/gtcp/pkg/tcpclient"
	"os"
	"strconv"
	"time"
)

func main() {
	clientLogger := logger.NewLogger(logger.LevelDebug)

	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "Specify a host, port, and message")
		os.Exit(1)
	}

	host := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Port must be an integer")
		os.Exit(1)
	}
	address := fmt.Sprintf("%s:%d", host, port)
	message := os.Args[3]

	for {
		recvd, err := tcpclient.Send(address, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending: %s", err)
		} else {
			clientLogger.Log(
				logger.LevelInfo,
				fmt.Sprintf(
					"Size: %d, IsCompressed: %v, Content: %s",
					recvd.GetSize(),
					recvd.IsCompressed(),
					recvd.GetContent(),
				),
			)
		}

		time.Sleep(1 * time.Second)
	}
}
