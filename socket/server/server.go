package main

import (
	"fmt"
	"log"

	"github.com/chinx/mohist/socket"
	"github.com/chinx/mohist/internal"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	server := "127.0.0.1:9111"

	socket.ProtocolHeader("mohist")
	if err := socket.ListenAndServe(server, handleConnection); err != nil {
		log.Println(fmt.Sprintf("Sever listen on %s is error: %s", server, err))
	}
}

func handleConnection(conn socket.Connect) {
	for {
		if msg, ok := <-conn.Receive(); ok {
			log.Println(fmt.Sprintf("Receive on service: %s msg: %s", conn.Addr(), internal.BytesString(msg)))
			conn.Send(msg)
		}
	}
}
