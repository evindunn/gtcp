package tcpmessage

import (
	"encoding/binary"
	"fmt"
	"gotest.tools/assert"
	"math/rand"
	"testing"
	"time"
)

// Latin unicode characters
const unicodeMin = 33
const unicodeMax = 95
const randStrSz = 1024

func getRandStr(strLen int) string {
	rand.Seed(time.Now().Unix())
	randRunes := make([]rune, strLen)
	for i:= range randRunes {
		randRunes[i] = rune(rand.Intn(unicodeMax-unicodeMin) + unicodeMin)
	}
	return string(randRunes)
}

func TestNewMessage(t *testing.T) {
	content := getRandStr(randStrSz)
	msg := NewMessage(content, false)

	assert.Assert(t, msg.isCompressed == 0, "NewMessage() stores incorrect value for isCompressed")
	msg = NewMessage(content, true)
	assert.Assert(t, msg.isCompressed == 1, "NewMessage() stores incorrect value for isCompressed")

	for i, val := range msg.content {
		assert.Assert(
			t,
			val == content[i],
			fmt.Sprintf("NewMessage() stores incorrect content byte at position %d", i))
	}
}

func TestMessage_GetSize(t *testing.T) {
	content := getRandStr(randStrSz)
	msg := NewMessage(content, false)

	assert.Assert(t, msg.GetSize() == len(msg.content), "Message.GetSize() returns incorrect value")
}

func TestMessage_IsCompressed(t *testing.T) {
	content := getRandStr(randStrSz)

	msg := NewMessage(content, false)
	assert.Assert(t, !msg.IsCompressed(), "Message.IsCompressed() returns incorrect value")

	msg.Compress()
	assert.Assert(t, msg.IsCompressed(), "Message.IsCompressed() returns incorrect value")

	msg.Decompress()
	assert.Assert(t, !msg.IsCompressed(), "Message.IsCompressed() returns incorrect value")
}

func TestMessage_GetContent(t *testing.T) {
	content := getRandStr(randStrSz)
	msg := NewMessage(content, false)

	msgContent := msg.GetContent()

	for i, val := range msgContent {
		assert.Assert(
			t,
			val == content[i],
			fmt.Sprintf("Message.GetContent() stores incorrect content byte at position %d", i))
	}
}

func TestMessage_ToBytes(t *testing.T) {
	content := getRandStr(randStrSz)
	msg := NewMessage(content, false)
	msgBytes := msg.ToBytes()

	msgSize := int(binary.LittleEndian.Uint64(msgBytes[0:8]))
	msgCompressed := int(msgBytes[8])
	msgContent := msgBytes[9:]

	assert.Assert(t, msgSize == len(msg.content), "Message.ToBytes() stores the wrong value for content size")
	assert.Assert(t, msgCompressed == msg.isCompressed, "Message.ToBytes() stores the wrong value for whether compressed")

	for i, val := range msgContent {
		assert.Assert(
			t,
			val == msg.content[i],
			fmt.Sprintf("Message.GetContent() stores incorrect content byte at position %d", i))
	}
}

func TestMessage_CompressDecompress(t *testing.T) {
	content := getRandStr(randStrSz)
	msg := NewMessage(content, false)
	szBeforeComp := msg.GetSize()
	msg.Compress()
	szAfterComp := msg.GetSize()
	msg.Decompress()
	msgContent := msg.GetContent()

	assert.Assert(t, szBeforeComp > szAfterComp, "Message.Compress() makes Message.content larger")

	for i, val := range msgContent {
		assert.Assert(
			t,
			val == content[i],
			fmt.Sprintf("Message.Compress() corrupts content byte at position %d", i))
	}
}