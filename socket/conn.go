package socket

import (
	"log"
	"net"
	"errors"
)

var (
	protocol Protocol
	noProtocolErr = errors.New("Please set protocal header first!")
)

type Handle func(Connect)

type Connect interface {
	Close() error
	Addr() string
	Receive() <-chan []byte
	Send(msg []byte) error
}

type connect struct {
	net.Conn
	User        string
	tmpBuffer   []byte
	readBuffer  []byte
	receiveChan chan []byte
}

func ProtocolHeader(header string) {
	protocol = NewProtocol(header)
}

func newConnect(conn net.Conn) *connect {
	c := &connect{
		Conn:        conn,
		tmpBuffer:   make([]byte, 0, 4096),
		readBuffer:  make([]byte, 2048),
		receiveChan: make(chan []byte, 2048),
	}
	go c.run()
	return c
}

func (c *connect) run() {
	for {
		n, err := c.Read(c.readBuffer)
		if err != nil {
			log.Println(c.RemoteAddr().String(), " connection error: ", err)
			return
		}

		c.tmpBuffer = protocol.Unpack(append(c.tmpBuffer, c.readBuffer[:n]...), c.receiveChan)
	}

	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

func (c *connect) Send(msg []byte) error {
	_, err := c.Write(protocol.Packet(msg))
	return err
}

func (c *connect) Receive() <-chan []byte {
	return c.receiveChan
}

func (c *connect) Addr() string {
	return c.RemoteAddr().String()
}

func ConnectTo(addr string, handler Handle) error {
	if protocol == nil{
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
	if protocol == nil{
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
