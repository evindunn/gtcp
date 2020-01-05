package tcpServer

import (
	"github.com/evindunn/gtcp/pkg/tcpMessage"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
)

type Server struct {
	listener *net.TCPListener
	connectionQueue chan *net.Conn
	signalQueue chan os.Signal
}

func NewServer(port int) (*Server, error) {
	addr := net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
		Zone: "",
	}
	listener, err := net.ListenTCP("tcp4", &addr)
	if err != nil {
		return nil, err
	}

	srv := Server{
		listener,
		make(chan *net.Conn),
		make(chan os.Signal, 1),
	}

	signal.Notify(srv.signalQueue, os.Interrupt)
	return &srv, nil
}

func handleMessage(c *net.Conn, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	connection := *c
	msg, err := tcpMessage.MessageFromConnection(c)
	if err != nil {
		return err
	}
	log.Printf("[%s] %s", connection.RemoteAddr().String(), string(msg.GetContent()))

	_, err = connection.Write([]byte("PONG"))
	if err != nil {
		return err
	}

	err = connection.Close()
	if err != nil {
		return err
	}
	log.Printf("Closed connection with %s\n", connection.RemoteAddr().String())

	return nil
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
			go handleMessage(c, &wg)
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
