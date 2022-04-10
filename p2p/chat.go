package p2p

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/guiwoo/gocoin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	// fmt.Printf("WebSocekt is connected %s",conn)
}
