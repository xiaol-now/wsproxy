package validate

import (
	"context"
	"wsproxy/middleware"
)

type Validator interface {
	Validate() error
}

func Middleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if v, ok := req.(Validator); ok {
				if err := v.Validate(); err != nil {
					return nil, err
				}
			}
			return handler(ctx, req)
		}
	}
}
