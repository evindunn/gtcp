# GoLang TCP

<a href="https://github.com/evindunn/gtcp/actions?query=workflow%3ABuild" target="_blank">
  <img type="image/svg" alt="gotest-status" src="https://github.com/evindunn/gtcp/workflows/Build/badge.svg">
</a>

<a href="https://github.com/evindunn/gtcp/actions?query=workflow%3ATest" target="_blank">
  <img type="image/svg" alt="gotest-status" src="https://github.com/evindunn/gtcp/workflows/Test/badge.svg">
</a>

<a href='https://coveralls.io/github/evindunn/gtcp?branch=master' target="_blank">
  <img src='https://coveralls.io/repos/github/evindunn/gtcp/badge.svg?branch=master&service=github&kill_cache=1' alt='Coverage Status' />
</a>

<a href="https://goreportcard.com/report/github.com/evindunn/gtcp" target="_blank">
  <img src="https://goreportcard.com/badge/github.com/evindunn/gtcp" alt="Go report card"/>
</a>

<a href="https://www.gnu.org/licenses/gpl-3.0" target="_blank">
    <img src="https://img.shields.io/badge/License-GPLv3-blue.svg" alt="License"/>
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
- [github.com/evindunn/gtcp/examples/tcpclient](./examples/tcpclient/main.go)
- [github.com/evindunn/gtcp/examples/tcpserver](./examples/tcpserver/main.go)

Handle connections using the interface defined in [ConnectionHandler.go](pkg/tcpserver/ConnectionHandler.go)
```go
type Handler struct {}

// This method makes the Handler type implement the ConnectionHandler interface
func (h *Handler) HandleConnection(c *net.Conn) {
    defer (*c).Close()
    log.Println("Connected!")
}
```

Then use the handler with the [Server](pkg/tcpserver/Server.go) struct
```go
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
