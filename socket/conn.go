package socket

import (
	"errors"
	"log"
	"net"
)

var (
	noProtocolErr = errors.New("Please set protocal header first!")
	protocol      Protocol
)

type Handle func(*Connect) bool

type Connect struct {
	conn        net.Conn
	Code        uint32
	User        string
	tmpBuffer   []byte
	readBuffer  []byte
	receiveChan chan []byte
}

func InitProtocol(kind, identifier string) {
	protocol = NewPacket(kind, identifier)
}

func newConnect(conn net.Conn) *Connect {
	c := &Connect{
		conn:        conn,
		tmpBuffer:   make([]byte, 0, 4096),
		readBuffer:  make([]byte, 2048),
		receiveChan: make(chan []byte, 2048),
	}
	go c.run()
	return c
}

func (c *Connect) run() {
	for {
		n, err := c.conn.Read(c.readBuffer)
		if err != nil {
			log.Println(c.conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		c.tmpBuffer = protocol.Unpack(append(c.tmpBuffer, c.readBuffer[:n]...), c.receiveChan)
	}

	if err := c.conn.Close(); err != nil {
		log.Println(err)
	}
}

func (c *Connect) Send(msg []byte) error {
	_, err := c.conn.Write(protocol.Packet(msg))
	return err
}

func (c *Connect) Receive() <-chan []byte {
	return c.receiveChan
}

func (c *Connect) Addr() string {
	return c.conn.RemoteAddr().String()
}

func ConnectTo(addr string, handler Handle) error {
	if protocol == nil {
		return noProtocolErr
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	handler(newConnect(conn))
	return nil
}

func ListenAndServe(addr string, handler Handle) error {
	if protocol == nil {
		return noProtocolErr
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handler(newConnect(conn))

	}
	return nil
}
