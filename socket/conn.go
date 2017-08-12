package socket

import "net"

type Connect interface {
	Run(string, func(net.Conn)) error
}

type client struct{}

type server struct{}

func NewClient() Connect { return &client{} }

func NewServer() Connect { return &server{} }

func (c *client) Run(addr string, handler func(net.Conn)) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	handler(conn)
	return nil
}

func (s *server) Run(addr string, handler func(net.Conn)) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handler(conn)
	}
	return nil
}
