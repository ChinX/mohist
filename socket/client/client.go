package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/chinx/mohist/internal"
	"github.com/chinx/mohist/socket"
)

type Msg struct {
	Meta    map[string]interface{} `json:"meta"`
	Content interface{}            `json:"content"`
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	server := "127.0.0.1:9111"

	socket.InitProtocol(socket.DefaultProtocol, "mohist")
	if err := socket.ConnectTo(server, receive); err != nil {
		log.Println(fmt.Sprintf("Client connect to %s is error: %s", server, err))
	}
}

func receive(conn *socket.Connect) bool {
	go send(conn)
	for {
		if msg, ok := <-conn.Receive(); ok {
			log.Println(fmt.Sprintf("Receive on service: %s msg: %s", conn.Addr(), internal.BytesString(msg)))
		}
	}
	return false
}

func send(conn *socket.Connect) {
	for i := 0; i < 100; i++ {
		session := GetSession()
		msg := &Msg{
			Meta: map[string]interface{}{
				"meta": "test",
				"id":   strconv.Itoa(i),
			},
			Content: Msg{
				Meta: map[string]interface{}{
					"author": "chinx",
				},
				Content: session,
			},
		}
		result, _ := json.Marshal(msg)
		if err := conn.Send(result); err != nil {
			log.Println(err)
		}
	}
}

func GetSession() string {
	gs1 := time.Now().Unix()
	gs2 := strconv.FormatInt(gs1, 10)
	return gs2
}
