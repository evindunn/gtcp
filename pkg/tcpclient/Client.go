package tcpclient

import (
	"fmt"
	"github.com/evindunn/gtcp/pkg/tcpmessage"
	"net"
	"os"
)

// Send a Message with content msgStr and return the Message response
func Send(addrStr string, msgStr string) (*tcpmessage.Message, error) {
	conn, err := net.Dial("tcp", addrStr)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return nil, err
	}

	msg := tcpmessage.NewMessage(msgStr, false)
	msgBytes := msg.ToBytes()
	bytesWritten, err := conn.Write(msgBytes)
	if err != nil {
		return nil, err
	}

	receivedMsg, err := tcpmessage.MessageFromConnection(&conn)
	if err != nil {
		return nil, err
	}

	actualBytes := len(msgBytes)
	if bytesWritten != actualBytes {
		fmt.Fprintf(
			os.Stderr,
			"Bytes written not equal to tcpmessage size: %d != %d", bytesWritten, actualBytes,
		)
	}

	return receivedMsg, nil
}
