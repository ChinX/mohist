package socket

import (
	"net"
	"log"
	"time"
)

func HeartBeating(conn net.Conn, readerChannel chan byte, timeout int)  {
	select {
	case <-readerChannel:
		log.Println(conn.RemoteAddr().String(), "get message, keeping heartbeating...")
		conn.SetDeadline(time.Now().Add(time.Duration(timeout)*time.Second))
		break
	case <-time.After(time.Duration(5)*time.Second):
		log.Println("It's really weird to get Nothing!!!")
		conn.Close()
	}
}

func GravelChannel(n []byte, mess chan byte){
	for _ , v := range n{
		mess <- v
	}
	close(mess)
}