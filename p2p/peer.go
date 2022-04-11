package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	key     string
	address string
	port    string
}

func AllPeers(p *peers) (result []string) {
	p.m.Lock()
	defer p.m.Unlock()
	for i, _ := range p.v {
		result = append(result, i)
	}
	return
}

func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.conn.Close()
	delete(Peers.v, p.key)
}

func (p *peer) read() {
	defer p.close()
	for {
		_, payload, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", payload)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		key:     key,
		address: address,
		port:    port,
	}
	go p.read()
	go p.write()

	Peers.v[key] = p
	return p
}
