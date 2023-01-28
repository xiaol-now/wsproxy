package wsproxy

import (
	"github.com/gorilla/websocket"
	"sync"
)

type writePack struct {
	typ     int
	content []byte
}

type Client struct {
	uuid string
	conn WebsocketConnInterface

	readCh  chan []byte
	writeCh chan writePack
	closeCh chan struct{}

	onceClose sync.Once
}

func NewClient(uuid string, conn *WSConn) *Client {
	c := &Client{
		uuid:    uuid,
		conn:    conn,
		readCh:  make(chan []byte, 3),
		writeCh: make(chan writePack, 3),
		closeCh: make(chan struct{}),
	}
	go c.readLoop()
	go c.writeLoop()
	return c
}

func (c *Client) Close() (err error) {
	c.onceClose.Do(func() {
		err = c.conn.Close()
		close(c.readCh)
		close(c.writeCh)
		close(c.closeCh)
	})
	return
}

func (c *Client) Wait() {
	<-c.closeCh
}

func (c *Client) IsHealth() bool {
	return c.conn.IsHealth()
}

func (c *Client) ReadMessage() []byte {
	return <-c.readCh
}

func (c *Client) SendMessage(message []byte) {
	select {
	case c.writeCh <- writePack{typ: websocket.TextMessage, content: message}:
		return
	default:
		// TODO; write已满，消息被拒绝
	}
}

func (c *Client) readLoop() {
	defer func(c *Client) { _ = c.Close() }(c)
	for {
		// TODO; messageType 待处理
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		select {
		case c.readCh <- payload:
		case <-c.closeCh:
			return
		}
	}

}

func (c *Client) writeLoop() {
	defer func() { _ = c.Close() }()
	for {
		select {
		case pack := <-c.writeCh:
			err := c.conn.WriteMessage(pack.typ, pack.content)
			if err != nil {
				return
			}
		case <-c.closeCh:
			return
		}
	}
}
