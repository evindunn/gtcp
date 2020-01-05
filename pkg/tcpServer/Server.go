package tcpServer

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
)

type Server struct {
	listener        net.TCPListener
	connectionQueue chan *net.Conn
	signalQueue     chan os.Signal
	handler         ConnectionHandler
}

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
		*listener,
		make(chan *net.Conn, 1024),
		make(chan os.Signal, 1),
		handler,
	}

	signal.Notify(srv.signalQueue, os.Interrupt)
	return &srv, nil
}

func (s *Server) handleMessage(c *net.Conn, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	s.handler.HandleConnection(c)
}

func (s *Server) queueConnections() {
	for {
		c, err := s.listener.Accept()
		if err == nil {
			s.connectionQueue <- &c
		}
	}
}

func (s *Server) Start() {
	log.Printf("Server listening on port %s...\n", strings.Split(s.listener.Addr().String(), ":")[1])

	var wg sync.WaitGroup

	go s.queueConnections()

	isRunning := true

	for {
		if !isRunning {
			break
		}

		select {
		case c := <- s.connectionQueue:
			go s.handleMessage(c, &wg)
			break
		default:
			break
		}

		select {
		case <-s.signalQueue:
			log.Println("Interrupt received, quitting")
			isRunning = false
			break
		default:
			break
		}
	}

	wg.Wait()

	close(s.connectionQueue)
	close(s.signalQueue)

	_ = s.listener.Close()
}
