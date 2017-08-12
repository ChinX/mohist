package main

import (
	"fmt"
	"log"
	"net"

	"github.com/chinx/mohist/internal"
	"github.com/chinx/mohist/socket"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	server := "127.0.0.1:9111"
	svc := socket.NewServer()
	if err := svc.Run(server, handleConnection); err != nil {
		log.Println(fmt.Sprintf("Sever listen on %s is error: %s", server, err))
	}
}

func handleConnection(conn net.Conn) {
	p := socket.NewProtocol("mohist")
	tmpBuffer := make([]byte, 0, 4096)
	buffer := make([]byte, 2048)
	receiver := make(chan []byte, 2048)
	go receive(receiver)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = p.Unpack(append(tmpBuffer, buffer[:n]...), receiver)
	}
	conn.Close()
}

func receive(receiver chan []byte) {
	for {
		msg, ok := <-receiver
		if !ok {
			break

		}
		log.Println("Receive msg string : ", internal.BytesString(msg))
	}
	log.Println("Receive msg over!")
}
