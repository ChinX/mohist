package socket

import (
	"bytes"

	"github.com/chinx/mohist/byteconv"
)

var (
	header    = byteconv.StrToBytes("Headers")
	hLen      = len(header)
	msgLen    = 4
	prefixLen = hLen + msgLen
	pool      = make([]byte, 0, 4096)
)

func Encode(msg []byte) []byte {
	return append(append(header, byteconv.IntToBytes(len(msg))...), msg...)
}

func Decode(buffer []byte) [][]byte {
	pool := append(pool, buffer...)
	buffers := make([][]byte, 0, 10)
	n := 0
	for i, l := 0, len(pool); i < l; i++ {
		if l < i+prefixLen {
			break
		}
		if bytes.Equal(pool[i:i+hLen], header) {
			ml := byteconv.BytesToInt(pool[i+hLen : i+prefixLen])
			if l < i+prefixLen+ml {
				break
			}
			buffers = append(buffers, pool[i+prefixLen:i+prefixLen+ml])
			n = i + prefixLen + ml - 1
		}
	}
	pool = pool[n:]
	return buffers
}
