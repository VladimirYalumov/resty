package resty

import (
	"context"
	"encoding/json"
	"github.com/VladimirYalumov/logger"
	"github.com/VladimirYalumov/resty/middleware"
	"github.com/VladimirYalumov/resty/requests"
	"github.com/VladimirYalumov/resty/responses"
	"github.com/VladimirYalumov/tracer"
	"github.com/rs/cors"
	"net/http"
)

var additionalMiddlewares []middleware.Middleware

type handler struct {
	*cors.Cors
	log *logger.Logger

	endpoints map[endpointKey]*endpoint
}

func NewHandler(log *logger.Logger, mm ...middleware.Middleware) *handler {
	additionalMiddlewares = make([]middleware.Middleware, len(mm)+1, len(mm)+1)
	j := 0
	for i := len(mm) - 1; i < 0; i-- {
		additionalMiddlewares[j] = mm[i]
		j++
	}
	additionalMiddlewares[j] = &middleware.RequestValidate{}

	return &handler{
		log:       log,
		endpoints: make(map[endpointKey]*endpoint),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer getDeferCatchPanic(h.log, w)

	ctx, span := tracer.StartSpan(context.Background(), r.URL.Path)
	span.Tag("method", r.Method)
	defer span.End()

	ctx = logger.ToContext(ctx, h.log.With("token", span.TraceId()))

	w.Header().Set("Content-Type", "application/json")

	e, ok := h.endpoints[endpointKey{r.URL.Path, r.Method}]
	if !ok || e == nil {
		logger.Warn(ctx, "unknown method", "method", r.Method, "path", r.URL.Path)
		w.WriteHeader(405)
		_ = json.NewEncoder(w).Encode(&responses.ErrorResponse{Message: "unknown method"})
		return
	}

	req := CheckAction(r, e.request, w)
	if req == nil {
		return
	}

	resp, httpCode := e.action(ctx, req)
	w.WriteHeader(httpCode)
	if err := resp.PrepareResponse(w); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		_, _ = w.Write([]byte{})
	}
	return
}

func (h *handler) Endpoint(path, method string, request requests.Request, action func(ctx context.Context, req requests.Request) (responses.Response, int), mm ...string) {
	e := &endpoint{
		action:      action,
		request:     request,
		middlewares: make(map[string]bool),
	}
	for _, m := range mm {
		e.middlewares[m] = true
	}
	e.middlewares[middleware.KeyRequestValidate] = true
	e.middlewares[middleware.KeyRequestInit] = true

	h.endpoints[endpointKey{path, method}] = e
}
