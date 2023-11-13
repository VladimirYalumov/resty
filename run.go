package resty

import (
	"context"
	"github.com/VladimirYalumov/logger"
	"github.com/rs/cors"
	"net/http"
	"resty/closer"
	"time"
)

const shutdownTimeout = 3 * time.Second

func RunServer[T any](ctx context.Context, h *handler[T], closerFns ...func(ctx context.Context) error) {
	c := &closer.Closer{}
	for _, closerFn := range closerFns {
		c.Add(closerFn)
	}

	go func() {
		if err := http.ListenAndServe(":8080", setCors(h)); err != nil {
			logger.Error(ctx, err, "serve")
		}
	}()
	logger.Info(ctx, "start server", "port", 8080)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(logger.ToContext(context.Background(), logger.FromContext(ctx)), shutdownTimeout)
	defer cancel()

	if err := c.Close(shutdownCtx); err != nil {
		logger.Error(ctx, err, "shutdown")
	}
	logger.Info(ctx, "stop")
}

func setCors[T any](handler *handler[T]) http.Handler {
	co := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "Client"},
		AllowCredentials: true,
	})

	return co.Handler(handler)
}
