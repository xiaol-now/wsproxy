package wsproxy

import (
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

type WebsocketConnInterface interface {
	IsHealth() bool
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type WSConn struct {
	lastConnect atomic.Int64
	conn        *websocket.Conn
}

func NewWSConn(conn *websocket.Conn) *WSConn {
	return &WSConn{conn: conn}
}

func (w *WSConn) IsHealth() bool {
	return time.Now().Unix()-w.lastConnect.Load() < 60
}

func (w *WSConn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = w.conn.ReadMessage()
	if err == nil {
		w.lastConnect.Store(time.Now().Unix())
	}
	return
}

func (w *WSConn) WriteMessage(messageType int, data []byte) error {
	err := w.conn.WriteMessage(messageType, data)
	if err == nil {
		w.lastConnect.Store(time.Now().Unix())
		return nil
	}
	return err
}

func (w *WSConn) Close() error {
	return w.conn.Close()
}
