package resty

import (
	"context"
	"fmt"
	"github.com/VladimirYalumov/logger"
	"github.com/VladimirYalumov/resty/closer"
	"net/http"
)

func RunServer(ctx context.Context, h *handler, opt Options, closerFns ...func(ctx context.Context) error) {
	c := &closer.Closer{}
	for _, closerFn := range closerFns {
		c.Add(closerFn)
	}

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", opt.Port), setCors(h)); err != nil {
			logger.Error(ctx, err, "serve")
		}
	}()
	logger.Info(ctx, "start server", "port", 8080)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(logger.ToContext(context.Background(), logger.FromContext(ctx)), opt.Timeout)
	defer cancel()

	if err := c.Close(shutdownCtx); err != nil {
		logger.Error(ctx, err, "shutdown")
	}
	logger.Info(ctx, "stop")
}
