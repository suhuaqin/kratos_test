package server

import (
	"net"
	"sync"
)

type SessionID string

type SessionMap struct {
	sessionMap map[SessionID]*session
	sync.RWMutex
}

type session struct {
	ID         SessionID
	Addr       *net.UDPAddr
	Send, Recv chan []byte
}

func newSessionMap() SessionMap {
	return SessionMap{
		sessionMap: make(map[SessionID]*session),
		RWMutex:    sync.RWMutex{},
	}
}

func newSession(id SessionID, addr *net.UDPAddr) *session {
	return &session{
		ID:   id,
		Addr: addr,
	}
}

func (sMap *SessionMap) addSession(id SessionID, addr *net.UDPAddr) (*session, bool) {
	sMap.RLocker()
	if sMap.sessionMap[id] != nil {
		sMap.RUnlock()
		return sMap.sessionMap[id], false
	}
	sMap.RUnlock()

	sMap.Lock()
	defer sMap.RWMutex.Unlock()
	sMap.sessionMap[id] = newSession(id, addr)
	return sMap.sessionMap[id], true
}
