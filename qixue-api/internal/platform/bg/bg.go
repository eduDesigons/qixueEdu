// Package bg runs work that must not delay the HTTP response.
// The task keeps request values but is not canceled when the client disconnects.
package bg

import (
	"context"
	"log/slog"
	"time"
)

const timeout = 8 * time.Second

func Go(parent context.Context, fn func(context.Context)) {
	if fn == nil {
		return
	}
	if parent == nil {
		parent = context.Background()
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("后台任务 panic", "recover", r)
			}
		}()
		ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), timeout)
		defer cancel()
		fn(ctx)
	}()
}
