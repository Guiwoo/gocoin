package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/guiwoo/gocoin/utils"
)

// this Upgrade is go routine!!
func chatUpgrade(rw http.ResponseWriter, r *http.Request) {
	var conns []*websocket.Conn
	var chatUpgrader = websocket.Upgrader{}
	chatUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := chatUpgrader.Upgrade(rw, r, nil)
	conns = append(conns, conn)
	utils.HandleErr(err)
	for {
		fmt.Println("Waiting meassage")
		_, p, err := conn.ReadMessage() //여기서 코드를 기다리네 ? 개쩌넹 ?
		if err != nil {
			conn.Close()
			break
		}
		for _, aConn := range conns {
			if aConn != conn {
				utils.HandleErr(aConn.WriteMessage(websocket.TextMessage, p))
			}
		}
	}
}
