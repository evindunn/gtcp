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
const HeaderSize = 9

/**
	TCP message struct
 */
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

/**
	Creates a new message from content
 */
func NewMessage(content string, isCompressed bool) Message {
	mContent := []byte(content)
	return Message{ content: mContent, isCompressed: boolToInt(isCompressed) }
}

/**
	|------- 8 bytes---------|------- 1 byte ---------|------- messageSize remaining bytes ---------|
	|----- messageSize ------|----- isCompressed -----|---------------- content -------------------|
 */
func (m *Message) ToBytes() []byte {
	msg := make([]byte, HeaderSize + m.GetSize())

	// 8 bytes for tcpmessage size
	binary.LittleEndian.PutUint64(msg, uint64(m.GetSize()))

	// One byte for whether compressed
	msg[8] = byte(m.isCompressed)

	// Remaining bytes are data
	msgContent := m.GetContent()
	for i := HeaderSize; i < HeaderSize + m.GetSize(); i++ {
		msg[i] = msgContent[i - HeaderSize]
	}

	return msg
}

/**
	Returns the content of the Message
 */
func (m *Message) GetContent() []byte {
	return m.content
}

/**
	Returns the length of the Message content
 */
func (m *Message) GetSize() int {
	return len(m.content)
}

/**
	Returns whether the Message is compressed
 */
func (m *Message) IsCompressed() bool {
	return m.isCompressed == 1
}

/**
	Compresses the Message content using zlib if the Message length exceeds compressionLimit
 */
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

/**
	Decompresses a compressed Message
 */
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

/**
	Reads a Message from a net.Conn
 */
func MessageFromConnection(c *net.Conn) (*Message, error)  {
	connReader := bufio.NewReader(*c)

	msgSizeBytes := make([]byte, 8)
	var isCompressedByte byte
	var content []byte
	var err error

	msgSize := 0
	var isCompressed int

	for i := 0; i > -1; i++ {
		if i < HeaderSize- 1 {
			msgSizeBytes[i], err = connReader.ReadByte()
			if err != nil {
				return nil, err
			}

		} else if i == HeaderSize- 1 {
			msgSize = int(binary.LittleEndian.Uint64(msgSizeBytes))
			content = make([]byte, msgSize)

			isCompressedByte, err = connReader.ReadByte()
			if err != nil {
				return nil, err
			}
			isCompressed = int(isCompressedByte)

		} else if msgSize > 0 && i < msgSize + HeaderSize {
			content[i - HeaderSize], err = connReader.ReadByte()
			if err != nil {
				return nil, err
			}

		} else {
			break
		}
	}

	return &Message{
		content: content,
		isCompressed: isCompressed,
	}, nil
}