package resty

import (
	"context"
	"encoding/json"
	"github.com/VladimirYalumov/logger"
	"github.com/VladimirYalumov/tracer"
	"net/http"
	"resty/action"
	"resty/middleware"
	"resty/requests"
	"resty/responses"
)

var additionalMiddlewares []middleware.Middleware

func Init(mm ...middleware.Middleware) {
	additionalMiddlewares = make([]middleware.Middleware, len(mm)+2)
	for i := len(mm) - 1; i != 0; i-- {
		additionalMiddlewares = append(additionalMiddlewares, mm[i])
	}
	additionalMiddlewares = append(additionalMiddlewares, &middleware.RequestValidate{})
}

type Handler[T any] struct {
	*cors.Cors
	log *logger.Logger

	endpoints map[endpointKey]*endpoint[T]
	data      *T
}

func (h *Handler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer getDeferCatchPanic(h.log, w)

	ctx, span := tracer.StartSpan(context.Background(), r.URL.Path)
	span.Tag("method", r.Method)
	defer span.End()

	ctx = logger.ToContext(ctx, h.log.With("token", span.TraceId()))

	w.Header().Set("Content-Type", "application/json")

	endpoint := action.GetEndpoint(r.Method, r.URL.Path)
	if endpoint == nil {
		logger.Warn(ctx, "unknown method", "method", r.Method, "path", r.URL.Path)
		w.WriteHeader(405)
		_ = json.NewEncoder(w).Encode(&responses.ErrorResponse{Message: "unknown method"})
		return
	}

	req := CheckAction(r, endpoint.Request(), w)
	if req == nil {
		return
	}

	endpoint.Action(ctx, req, w)
	return
}

func (h *Handler[T]) Endpoint(
	method,
	path string,
	request requests.Request, action func(ctx context.Context, req requests.Request, w http.ResponseWriter),
	mm ...string,
) {
	key := endpointKey{path, method}
	h.endpoints[key] = &endpoint[T]{method: method, Action: action, request: request, data: h.data}
	for _, m := range mm {
		h.endpoints[key].middlewares[m] = true
	}
	h.endpoints[key].middlewares[middleware.KeyRequestValidate] = true
	h.endpoints[key].middlewares[middleware.KeyRequestValidate] = true
}

func (h *Handler[T]) SetAdditionalData(data *T) {
	h.data = data
}
