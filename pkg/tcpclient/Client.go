package tcpclient

import (
	"fmt"
	"github.com/evindunn/gtcp/pkg/tcpmessage"
	"net"
)

// Send a Message with content msgStr
func Send(addrStr string, msgStr string) error {
	conn, err := net.Dial("tcp", addrStr)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return err
	}

	msg := tcpmessage.NewMessage(msgStr, false)
	msgBytes := msg.ToBytes()
	bytesWritten, err := conn.Write(msgBytes)
	if err != nil {
		return err
	}

	actualBytes := len(msgBytes)
	if bytesWritten != actualBytes {
		return fmt.Errorf("bytes written not equal to tcpmessage size: %d != %d", bytesWritten, actualBytes)
	}

	return nil
}
