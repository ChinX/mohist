package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/chinx/mohist/socket"
)

type Msg struct {
	Meta    map[string]interface{} `json:"meta"`
	Content interface{}            `json:"content"`
}

func main() {
	server := "127.0.0.1:9111"

	cli := socket.NewClient()
	if err := cli.Run(server, send); err != nil {
		log.Println(fmt.Sprintf("Client connect to %s is error: %s", server, err))
	}
}

func send(conn net.Conn) {
	protocol := socket.NewProtocol("mohist")
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
		log.Println(string(result))
		conn.Write(protocol.Packet(result))
	}
	fmt.Println("send over")
	select {
	case <-time.After(time.Duration(20) * time.Second):
		conn.Close()
	}
}

func GetSession() string {
	gs1 := time.Now().Unix()
	gs2 := strconv.FormatInt(gs1, 10)
	return gs2
}
