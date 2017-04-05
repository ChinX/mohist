package socket

import (
	"log"
	"net"
)

var listen net.Listener

type socket struct {
	conn net.Conn
}

func Run(laddr string) {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		panic(err)
	}
	listen = l
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		//conn.
		log.Println(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//接收解包
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = Depacket(append(tmpBuffer, buffer[:n]...))
		if len(tmpBuffer) > 0 {
			readerChannel <- tmpBuffer
		}
	}
	conn.Close()
}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			log.Println(string(data))
		}
	}
}
