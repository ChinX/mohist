package socket

import (
	"log"
	"net"
)

var protocol = NewProtocol("mohist")

type Transmitter struct {
	net.Conn
}

func (c *Transmitter) Send(msg []byte) {
	c.Write(protocol.Packet(msg))
}

func (c *Transmitter) Receive(receiver chan []byte) {
	tmpBuffer := make([]byte, 0, 4096)
	buffer := make([]byte, 2048)
	for {
		n, err := c.Conn.Read(buffer)
		if err != nil {

			log.Println(c.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), receiver)
	}
}
