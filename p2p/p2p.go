package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/guiwoo/gocoin/utils"
)

var upgrader = websocket.Upgrader{}

// this Upgrade is go routine!!
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Port:3000 will upgrade the request from :4000

	openPort := r.URL.Query().Get("openPort")
	ip := utils.Spliter(r.RemoteAddr, ":", 0)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)
}

func AddPeer(address, port, openPort string) {
	// @from :4000 => is requesting an upgrade from the port :3000
	// nil part on maybe token or somthing like authentications. cookie token etc...
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleErr(err)
	initPeer(conn, address, port)
}
