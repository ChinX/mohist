package socket

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func send(conn net.Conn) {
	for i := 0; i < 100; i++ {
		session := GetSession()
		words := "{\"ID\":" + strconv.Itoa(i) + "\",\"Session\":" + session + "2015073109532345\",\"Meta\":\"golang\",\"Content\":\"message\"}"
		conn.Write(Enpacket([]byte(words)))
	}
	fmt.Println("send over")
}

func GetSession() string {
	gs1 := time.Now().Unix()
	gs2 := strconv.FormatInt(gs1, 10)
	return gs2
}

func Connect(laddr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", laddr)
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
