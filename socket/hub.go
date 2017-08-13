package socket

import "sync"

var conHub = &ConnHub{hash: make(map[string]*connect)}

type ConnHub struct {
	sync.RWMutex
	hash map[string]*connect
}

func Register(connect *connect) {
	conHub.Lock()
	conHub.hash[connect.User] = connect
	conHub.Unlock()
}

func Deregister(connect *connect) {
	conHub.Lock()
	delete(conHub.hash, connect.User)
	conHub.Unlock()
}

func InvokeEach(handler Handle) {
	conHub.RLock()
	for _, conn := range conHub.hash {
		go handler(conn)
	}
	conHub.RUnlock()
}

func Invoke(handler Handle, users ...string) {
	conHub.RLock()
	for user, conn := range conHub.hash {
		for i, l := 0, len(users); i < l; i++ {
			if user == users[i] {
				go handler(conn)
			}
		}
	}
	conHub.RUnlock()
}
