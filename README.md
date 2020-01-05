# GoLang TCP

<a href="https://github.com/evindunn/gtcp/actions?query=workflow%3ABuild">
  <img type="image/svg" alt="gotest-status" src="https://github.com/evindunn/gtcp/workflows/Build/badge.svg">
</a>

<a href="https://github.com/evindunn/gtcp/actions?query=workflow%3ATest">
  <img type="image/svg" alt="gotest-status" src="https://github.com/evindunn/gtcp/workflows/Test/badge.svg">
</a>

<a href='https://coveralls.io/github/evindunn/gtcp?branch=master'>
  <img src='https://coveralls.io/repos/github/evindunn/gtcp/badge.svg?branch=master&service=github' alt='Coverage Status' />
</a>

<a href="https://goreportcard.com/report/github.com/evindunn/gtcp">
  <img src="https://goreportcard.com/badge/github.com/evindunn/gtcp" alt="Go report card"/>
</a>

Simple TCP message passing in Go

#### API
- github.com/evindunn/gtcp/pkg/tcpclient
  - Send(addrStr string, msgStr string) error
- github.com/evindunn/gtcp/pkg/tcpserver
  - NewServer(port int, handler ConnectionHandler) (*Server, error)
  - Server.Start()
- github.com/evindunn/gtcp/pkg/tcpmessage
  - NewMessage(content string, isCompressed bool) Message
  - MessageFromConnection(c *net.Conn) (*Message, error)
  - Message.ToBytes() []byte
  - Message.GetContent() []byte
  - Message.GetSize() int
  - Message.IsCompressed() bool
  - Message.Compress() error
  - Message.Decompress() error
  
#### Examples
- github.com/evindunn/gtcp/cmd/tcpclient
- github.com/evindunn/gtcp/cmd/tcpserver

Handle connections using the interface defined in [ConnectionHandler.go](pkg/tcpserver/ConnectionHandler.go)
```text
type Handler struct {}

// This method makes the Handler type implement the ConnectionHandler interface
func (h *Handler) HandleConnection(c *net.Conn) {
    log.Println("Connected!")
    (*c).Close()
}
```

Then use the handler with the [Server](pkg/tcpserver/Server.go) struct
```text
var h Handler
srv, err := tcpServer.NewServer(8080, &h)

if err != nil {
    fmt.Fprintf(os.Stderr, "An error occurred creating the server: %s\n", err)
    os.Exit(1)
}

srv.Start()
```

The [Message](pkg/tcpserver/Server.go) class defines a simple protocol for passing messages in TCP
```text
|------- 8 bytes---------|------- 1 byte ---------|------- Remaining bytes ---------|
|----- messageSize ------|----- isCompressed -----|---------- content --------------|
```

Receiving a message server-side ([example](cmd/tcpserver/main.go)):
```text
type Handler struct {}

func (h *Handler) HandleConnection(c *net.Conn) {
	msg, err := tcpMessage.MessageFromConnection(c)
	if err != nil {
		log.Printf("[%s] Error parsing connection: %s\n", (*c).RemoteAddr().String(), err)
	} else {
        log.Printf(
            "[%s] Size: %d, Compressed: %v, Content: %s", connection.RemoteAddr().String(),
            msg.GetSize(),
            msg.IsCompressed(),
            string(msg.GetContent()))
    }

    (*c).Close()
}

func main() {
    var h Handler
    port := 8080

    srv, _ := tcpServer.NewServer(port, &h)
    srv.Start()
}
```

Sending a Message client-side ([example in tcpClient.Send()](pkg/tcpclient/Client.go)):
```text
addrStr := "127.0.0.1:8080"
conn, _ := net.Dial("tcp", addrStr)
defer conn.Close()

msg := tcpMessage.NewMessage(msgStr, false)
msgBytes := msg.ToBytes()

bytesWritten, err := conn.Write(msgBytes)
if err != nil {
    return err
}
```