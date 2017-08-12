package socket

import (
	"bytes"

	"github.com/chinx/mohist/internal"
)

const (
	dataLen = 4
)

//通讯协议处理，主要处理封包和解包的过程
type Protocol interface {
	Packet([]byte) []byte
	Unpack([]byte, chan<- []byte) []byte
}

type protocal struct {
	header    []byte
	headerLen int
	prefixLen int
}

func NewProtocol(packetHeader string) Protocol {
	header := internal.StringBytes(packetHeader)
	hLen := len(header)
	return &protocal{
		header:    header,
		headerLen: hLen,
		prefixLen: hLen + dataLen,
	}
}

// 封包
func (p *protocal) Packet(msg []byte) []byte {
	msgLen := len(msg)
	buffer := make([]byte, int(p.prefixLen)+msgLen)
	copy(buffer[:p.headerLen], p.header)
	copy(buffer[p.headerLen:p.prefixLen], internal.Uint32ToBytes(uint32(msgLen)))
	copy(buffer[p.prefixLen:], msg)
	return buffer
}

// 解包
func (p *protocal) Unpack(buffer []byte, receiver chan<- []byte) []byte {
	n := 0
	for i, l := 0, len(buffer); i < l; i++ {
		starting := i + p.prefixLen
		if l < starting {
			break
		}
		if bytes.Equal(buffer[i:i+p.headerLen], p.header) {
			ending := starting + int(internal.BytesToUint32(buffer[i+p.headerLen:starting]))
			if l < ending {
				break
			}
			receiver <- buffer[starting:ending]
			n = ending - 1
		}
	}
	return buffer[n:]
}
