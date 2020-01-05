package tcpserver

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Server wraps net.TCPListener to simplify multi-threaded message passing
type Server struct {
	listener    *net.TCPListener
	signalQueue chan os.Signal
	handler     ConnectionHandler
}

// NewServer creates a new server that runs on port and handles connections using handler
func NewServer(port int, handler ConnectionHandler) (*Server, error) {
	addr := net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
		Zone: "",
	}
	listener, err := net.ListenTCP("tcp4", &addr)
	if err != nil {
		return nil, err
	}

	// TODO: Max connections
	srv := Server{
		listener,
		make(chan os.Signal, 1),
		handler,
	}

	signal.Notify(srv.signalQueue, os.Interrupt)
	return &srv, nil
}

func (s *Server) handleConnection(c *net.Conn, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	s.handler.HandleConnection(c)
}

func (s *Server) watchInterrupts() {
	<-s.signalQueue
	s.Stop()
}

// Stop sets the deadling for Server.listener to time.Now()
func (s *Server) Stop() {
	s.listener.SetDeadline(time.Now())
}

/*
Start starts the server
A goroutine called queueConnections() does just that
Connections in the queue are dispatched to goroutines that call handleMessage()
Listens for the interrupt signal and terminates gracefully when encountered
*/
func (s *Server) Start() {
	defer close(s.signalQueue)
	defer s.listener.Close()

	go s.watchInterrupts()
	var wg sync.WaitGroup

	for {
		c, err := s.listener.Accept()
		if err == nil {
			go s.handleConnection(&c, &wg)
		} else {
			break
		}
	}

	wg.Wait()
}
