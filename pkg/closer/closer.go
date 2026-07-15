package closer

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type closeFn struct {
	name string
	fn   func(context.Context) error
}

type closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []closeFn
}

var globalCloser = newCloser()

func newCloser() *closer {
	return &closer{}
}

func Add(name string, f func(context.Context) error) {
	globalCloser.Add(name, f)
}

func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

func (c *closer) Add(name string, f func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, closeFn{name: name, fn: f})
}

func (c *closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			return
		}

		slog.Info("starting graceful shutdown", "count", len(funcs))

		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]

			start := time.Now()
			slog.Info("закрываем ресурс", "name", f.name)

			if err := f.fn(ctx); err != nil {
				slog.Error("error while closing resource ", "name", f.name, "error", err, "duration", time.Since(start))

				if result == nil {
					result = err
				}
			} else {
				slog.Info("closed", "name", f.name, "duration", time.Since(start))
			}
		}

		slog.Info("graceful shutdown finished", "count", len(funcs))
	})

	return result
}
