package main

import (
	"fmt"
	"github.com/chinx/mohist/socket"
	"k8s.io/kubernetes/pkg/util/json"
	"net"
	"os"
	"strconv"
	"time"
)

type Msg struct {
	Meta    map[string]interface{} `json:"meta"`
	Content interface{}            `json:"content"`
}

func send(conn net.Conn) {
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
		conn.Write(socket.Encode(result))
	}
	fmt.Println("send over")
	select {
	case <-time.After(time.Duration(10) * time.Second):
		conn.Close()
	}
}

func GetSession() string {
	gs1 := time.Now().Unix()
	gs2 := strconv.FormatInt(gs1, 10)
	return gs2
}

func main() {
	server := "127.0.0.1:9111"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	send(conn)

}
