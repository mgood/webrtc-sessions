package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/progrium/webrtc-sessions/sfu"
	"github.com/progrium/webrtc-sessions/web"
	"tractor.dev/toolkit-go/engine"
)

func main() {
	engine.Run(Main{})
}

type Main struct{}

func (m *Main) Serve(ctx context.Context) {

	session := sfu.NewSession()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		peer, err := session.AddPeer(conn)
		if err != nil {
			log.Print("peer:", err)
			return
		}
		peer.HandleSignals()
	})

	http.Handle("/", http.FileServer(http.FS(web.Dir)))

	log.Println("running on http://localhost:8088 ...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
