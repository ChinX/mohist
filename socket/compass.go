package socket

import "fmt"

var (
	compMap = make(map[uint32]Handle)
	notFound = notFoundHandler
)


func Add(code uint32, handles ...Handle)  {
	if _, ok := compMap[code]; ok {
		panic(fmt.Sprintf("Registry code:%d is exist!", code))
	}
	compMap[code] = handlersChain(handles)
}

func Match(code uint32) Handle {
	handler, ok := compMap[code]
	if !ok {
		handler = notFound
	}
	return handler
}

func NotFoundHandler(handler Handle)  {
	notFound = handler
}

func notFoundHandler(conn *Connect) bool {
	conn.Send(Integrate(conn.Code, []byte("Not found")))
	return true
}

func handlersChain(handlers []Handle) Handle {
	nHandlers := make([]Handle, 0, 10)
	l := 0
	for i := 0; i < len(handlers); i++ {
		if handlers[i] != nil {
			nHandlers = append(nHandlers, handlers[i])
			l++
		}
	}
	if l == 0 {
		return nil
	}
	return func(conn *Connect) bool {
		length := len(handlers)
		for i := 0; i < length; i++ {
			if handlers[i](conn){
				return true
			}
		}
		return true
	}
}
