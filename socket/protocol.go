package socket

import (
	"bytes"

	"github.com/chinx/mohist/byteconv"
)

var (
	header    []byte = byteconv.StrToBytes("Headers")
	hLen             = len(header)
	msgLen           = 4
	prefixLen        = hLen + msgLen
)

func SetHandler(h string) {
	header = byteconv.StrToBytes(h)
	hLen = len(header)
	prefixLen = hLen + msgLen
}

func Enpacket(msg []byte) []byte {
	return append( append(header, byteconv.IntToBytes(len(msg))...), msg...)
}

func Depacket(buffer []byte) []byte {
	i, l := 0, len(buffer)
	data := make([]byte, 2048)
	for ; i < l; i++ {
		if l < i+prefixLen {
			break
		}
		if bytes.Equal(buffer[i:i+hLen], header) {
			ml := byteconv.BytesToInt(buffer[i+hLen : i+prefixLen])
			if l < i+prefixLen+ml {
				break
			}
			data = buffer[i+prefixLen : i+prefixLen+ml]
		}
	}
	if i == l {
		return make([]byte, 0)
	}
	return data
}
