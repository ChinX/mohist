package socket

import "sync"

var conHub = &ConnHub{hash: make(map[string]*Connect)}

type ConnHub struct {
	sync.RWMutex
	hash map[string]*Connect
}

func Register(connect *Connect) {
	conHub.RLock()
	if _, ok := conHub.hash[connect.User]; ok {
		//conn.Conn.Write()
		//closed
	}
	conHub.RUnlock()
	conHub.Lock()
	conHub.hash[connect.User] = connect
	conHub.Unlock()
}

func Deregister(connect *Connect) {
	conHub.RLock()
	if _, ok := conHub.hash[connect.User]; ok {
		conHub.RUnlock()
		conHub.Lock()
		delete(conHub.hash, connect.User)
		conHub.Unlock()
		//conn.Conn.Write()
		//closed
	}else{
		conHub.RUnlock()
	}
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
