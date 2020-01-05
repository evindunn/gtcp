package tcpmessage

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"net"
)

const compressionLimit = 1024

/*
HeaderSize represents the number of bytes that make up the Message header,
which includes 8 bytes for content size and 1 byte for whether the Message
is compressed
*/
// which includes
const HeaderSize = 9

// Message represents a message sent over TCP
type Message struct {
	content      []byte
	isCompressed int
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// NewMessage creates a new message from content
func NewMessage(content string, isCompressed bool) Message {
	mContent := []byte(content)
	return Message{content: mContent, isCompressed: boolToInt(isCompressed)}
}

/*
ToBytes exports Message to the following format:
|------- 8 bytes---------|------- 1 byte ---------|------- messageSize remaining bytes ---------|
|----- messageSize ------|----- isCompressed -----|---------------- content --------------------|
*/
func (m *Message) ToBytes() []byte {
	msg := make([]byte, HeaderSize+m.GetSize())

	// 8 bytes for tcpmessage size
	binary.LittleEndian.PutUint64(msg, uint64(m.GetSize()))

	// One byte for whether compressed
	msg[8] = byte(m.isCompressed)

	// Remaining bytes are data
	msgContent := m.GetContent()
	for i := HeaderSize; i < HeaderSize+m.GetSize(); i++ {
		msg[i] = msgContent[i-HeaderSize]
	}

	return msg
}

// GetContent returns the content of the Message without header
func (m *Message) GetContent() []byte {
	return m.content
}

// GetSize returns the length of the Message content
func (m *Message) GetSize() int {
	return len(m.content)
}

// IsCompressed Returns whether the Message is compressed
func (m *Message) IsCompressed() bool {
	return m.isCompressed == 1
}

// Compress uses zlib to compress the Message if the Message length exceeds compressionLimit
func (m *Message) Compress() error {
	if !m.IsCompressed() && m.GetSize() >= compressionLimit {
		var contentBuf bytes.Buffer
		compressor := zlib.NewWriter(&contentBuf)
		_, err := compressor.Write(m.content)
		if err != nil {
			return err
		}
		err = compressor.Close()
		if err != nil {
			return err
		}
		m.content = contentBuf.Bytes()
		m.isCompressed = 1
	}

	return nil
}

// Decompress decompresses a compressed Message
func (m *Message) Decompress() error {
	if m.IsCompressed() {
		decompressor, err := zlib.NewReader(bytes.NewReader(m.content))
		if err != nil {
			return err
		}

		var mContent bytes.Buffer
		contentWriter := bufio.NewWriter(&mContent)
		_, err = io.Copy(contentWriter, decompressor)
		if err != nil {
			return err
		}

		m.content = mContent.Bytes()

		err = decompressor.Close()
		if err != nil {
			return err
		}
		m.isCompressed = 0
	}

	return nil
}

// MessageFromConnection reads a Message from a net.Conn
func MessageFromConnection(c *net.Conn) (*Message, error) {
	connReader := bufio.NewReader(*c)

	var err error
	var msgSize int
	var isCompressed int
	var content []byte

	// Read header
	header := make([]byte, HeaderSize)
	for i := 0; i < len(header); i++ {
		header[i], err = connReader.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	msgSize = int(binary.LittleEndian.Uint64(header[:HeaderSize-1]))
	isCompressed = int(header[HeaderSize-1])
	content = make([]byte, msgSize)

	// Read content
	for i := 0; i < len(content); i++ {
		content[i], err = connReader.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		content:      content,
		isCompressed: isCompressed,
	}, nil
}
