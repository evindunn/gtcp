package main

import (
	"fmt"
	"github.com/evindunn/gtcp/pkg/tcpclient"
	"os"
	"strconv"
	"sync"
)

func asyncSend(address string, message string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	err := tcpclient.Send(address, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func main() {
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
	message := os.Args[3]

	counter := 0
	var wg sync.WaitGroup
	for {
		if counter >= 10 {
			break
		}

		fmt.Printf("Sent %d\n", counter)
		asyncSend(fmt.Sprintf("%s:%d", host, port), fmt.Sprintf("%s%d", message, counter), &wg)

		counter += 1
	}

	wg.Wait()
}
