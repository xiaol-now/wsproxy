package wsproxy

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"wsproxy/middleware"
)

type WSProxy struct {
	upgrader    *websocket.Upgrader
	middlewares []middleware.Middleware
}

func NewWSProxy(upgrader *websocket.Upgrader) *WSProxy {
	return &WSProxy{
		upgrader:    upgrader,
		middlewares: make([]middleware.Middleware, 0),
	}
}

func (ws *WSProxy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !websocket.IsWebSocketUpgrade(request) {
		// TODO; 已经是ws协议的情况下如何处理
		http.Error(writer, "request error", 500)
		return
	}
	handler := middleware.Chain(ws.middlewares...)(func(ctx context.Context, req any) (any, error) {
		return req, ws.handler(ctx.(Context))
	})
	_, err := handler(NewWrapContext(writer, request), nil) // TODO; 传递的参数
	if err != nil {
		//TODO; 如何处理中间件返回的错误
		http.Error(writer, "request error", 500)
		return
	}
}

func (ws *WSProxy) handler(ctx Context) error {
	var responseHeader http.Header
	conn, err := ws.upgrader.Upgrade(ctx.Response(), ctx.Request(), responseHeader)
	if err != nil {
		return err
	}
	client := NewClient(uuid.New().String(), NewWSConn(conn))
	defer func() { _ = client.Close() }()

	for {
		message := client.ReadMessage()
		if message == nil {
			return nil
		}
		client.SendMessage(append([]byte("已经收到消息："), message...))
	}
}

func (ws *WSProxy) Use(m ...middleware.Middleware) {
	ws.middlewares = append(ws.middlewares, m...)
}
