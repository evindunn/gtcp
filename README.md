# GoLang TCP

<a href="https://github.com/evindunn/Tcp/actions?query=workflow%3ABuild">
  <img type="image/svg" alt="gotest-status" src="https://github.com/evindunn/gtcp/workflows/Build/badge.svg">
</a>

<a href="https://github.com/evindunn/Tcp/actions?query=workflow%3ATest">
  <img type="image/svg" alt="gotest-status" src="https://github.com/evindunn/gtcp/workflows/Test/badge.svg">
</a>

<a href='https://coveralls.io/github/evindunn/gtcp?branch=master'>
  <img src='https://coveralls.io/repos/github/evindunn/gtcp/badge.svg?branch=master&service=github' alt='Coverage Status' />
</a>

A message-passer built on TCP written in GO
- zlib compression for larger messages
- Utilities for converting the Tcp.Message to and from raw bytes
- Utilities for reading/writing Tcp.Message to/from net.Conn

Handle connections using the interface defined in [ConnectionHandler.go](./pkg/tcpServer/ConnectionHandler.go)
```text
type Handler struct {}

// This method makes the Handler type implement the ConnectionHandler interface
func (h *Handler) HandleConnection(c *net.Conn) {
    log.Println("Connected!")
    (*c).Close()
}
```

Then use the handler with the [Server](./pkg/tcpServer/Server.go) struct
```text
var h Handler
srv, err := tcpServer.NewServer(8080, &h)

if err != nil {
    fmt.Fprintf(os.Stderr, "An error occurred creating the server: %s\n", err)
    os.Exit(1)
}

srv.Start()
```

The [Message](./pkg/tcpServer/Server.go) class defines a simple protocol for passing messages in TCP
```text
|------- 8 bytes---------|------- 1 byte ---------|------- Remaining bytes ---------|
|----- messageSize ------|----- isCompressed -----|---------- content --------------|
```

Receiving a message server-side ([example](./cmd/tcpServer/main.go)):
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

Sending a Message client-side ([example](./pkg/tcpClient/Client.go)):
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