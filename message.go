package tcpless

import (
	"encoding/binary"
	"io"
)

// Message standard tcp message
type Message struct {
	// route name length
	handlerNameLength uint8
	// route name
	handlerName string
	// body length
	streamLength uint64
	// body data
	stream io.Reader
}

// CreateMessage create message func
// handler - name of handler
// data - body
func CreateMessage(handler string, data []byte) []byte {
	if len(handler) > 255 {
		return nil
	}
	dataLen := len(data)
	result := make([]byte, len(handler)+dataLen+9)
	result[0] = uint8(len(handler))
	copy(result[1:result[0]+1], handler)
	binary.BigEndian.PutUint64(result[result[0]+1:result[0]+9], uint64(dataLen))
	copy(result[result[0]+9:], data)
	return result
}
