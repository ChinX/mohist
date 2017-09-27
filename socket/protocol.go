package socket

import (
	"bytes"

	"github.com/chinx/mohist/internal"
)

const (
	DataLen         = 4
	DefaultProtocol = "default"
)

//通讯协议处理，主要处理封包和解包的过程
type Protocol interface {
	Packet([]byte) []byte
	Unpack([]byte, chan<- []byte) []byte
}

type packet struct {
	identifier    []byte
	identifierLen int
	prefixLen     int
}

func NewPacket(kind, identifier string) (ptl Protocol) {
	switch kind {
	case DefaultProtocol:
		fallthrough
	default:
		header := internal.StringBytes(identifier)
		hLen := len(header)
		ptl = &packet{
			identifier:    header,
			identifierLen: hLen,
			prefixLen:     hLen + DataLen,
		}
	}
	return ptl
}

// 封包
func (p *packet) Packet(msg []byte) []byte {
	msgLen := len(msg)
	buffer := make([]byte, p.prefixLen+msgLen)
	copy(buffer[:p.identifierLen], p.identifier)
	copy(buffer[p.identifierLen:p.prefixLen], internal.Uint32ToBytes(uint32(msgLen)))
	copy(buffer[p.prefixLen:], msg)
	return buffer
}

// 解包
func (p *packet) Unpack(buffer []byte, receiver chan<- []byte) []byte {
	n := 0
	for i, l := 0, len(buffer); i < l; i++ {
		starting := i + p.prefixLen
		if l < starting {
			break
		}
		if bytes.Equal(buffer[i:i+p.identifierLen], p.identifier) {
			ending := starting + int(internal.BytesToUint32(buffer[i+p.identifierLen:starting]))
			if l < ending {
				break
			}
			receiver <- buffer[starting:ending]
			n = ending - 1
		}
	}
	return buffer[n:]
}

func Integrate(code uint32, byteArr []byte) []byte {
	buffer := make([]byte, len(byteArr)+DataLen)
	copy(buffer[:DataLen], internal.Uint32ToBytes(code))
	copy(buffer[DataLen:], byteArr)
	return buffer
}

func Disintrgeate(byteArr []byte) (uint32, []byte) {
	return internal.BytesToUint32(byteArr[:DataLen]), byteArr[DataLen:]
}
