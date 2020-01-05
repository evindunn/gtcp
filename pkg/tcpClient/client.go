package tcpClient

import (
	"fmt"
	"github.com/evindunn/gtcp/pkg/tcpMessage"
	"net"
)

func Send(addrStr string, msgStr string) error {
	conn, err := net.Dial("tcp", addrStr)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return err
	}

	msg := tcpMessage.NewMessage(msgStr, false)
	msgBytes := msg.ToBytes()
	bytesWritten, err := conn.Write(msgBytes)
	if err != nil {
		return err
	}

	actualBytes := len(msgBytes)
	if bytesWritten != actualBytes {
		return fmt.Errorf("bytes written not equal to tcpMessage size: %d != %d", bytesWritten, actualBytes)
	}

	return nil
}