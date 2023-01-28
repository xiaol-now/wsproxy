package wsproxy

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
	"time"
	"wsproxy/middleware/trace"
	"wsproxy/middleware/validate"
)

func TestWSProxy_ServeHTTP(t *testing.T) {
	wsp := NewWSProxy(&websocket.Upgrader{
		HandshakeTimeout: time.Second * 30,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})
	wsp.Use(validate.Middleware(), trace.Middleware())

	http.Handle("/ws", wsp)
	log.Fatal(http.ListenAndServe(":10000", wsp))
}
