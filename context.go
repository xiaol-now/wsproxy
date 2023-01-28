package wsproxy

import (
	"context"
	"net/http"
)

type Context interface {
	context.Context
	Request() *http.Request
	Response() http.ResponseWriter
}

type WrapContext struct {
	context.Context
	request *http.Request
	writer  http.ResponseWriter
}

func NewWrapContext(writer http.ResponseWriter, request *http.Request) *WrapContext {
	return &WrapContext{
		Context: context.Background(),
		request: request,
		writer:  writer,
	}
}

func (w *WrapContext) Request() *http.Request {
	return w.request
}

func (w *WrapContext) Response() http.ResponseWriter {
	return w.writer
}
