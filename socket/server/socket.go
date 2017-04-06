package main

import (
	"log"
	"net"

	"github.com/chinx/mohist/socket"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	server := "127.0.0.1:9111"
	l, err := net.Listen("tcp", server)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		//conn.
		log.Println(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffers := socket.Decode(buffer[:n])
		for i := 0; i < len(tmpBuffers); i++ {
			log.Println("receive data string:", string(tmpBuffers[i]))
		}

	}
	conn.Close()
}
