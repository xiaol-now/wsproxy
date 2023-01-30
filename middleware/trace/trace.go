package trace

import "wsproxy/middleware"

func Middleware() middleware.Middleware {
	// TODO; 链路追踪
	return func(next middleware.Handler) middleware.Handler {
		return next
	}
}
