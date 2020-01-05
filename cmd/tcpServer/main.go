package main

import (
	"fmt"
	"github.com/evindunn/gtcp/pkg/tcpServer"
	"os"
	"strconv"
)

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

	srv, err := tcpServer.NewServer(port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred creating the server: %s\n", err)
		os.Exit(1)
	}

	srv.Start()
}