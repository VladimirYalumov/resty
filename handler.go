package resty

import (
	"context"
	"encoding/json"
	"github.com/VladimirYalumov/logger"
	"github.com/VladimirYalumov/resty/middleware"
	"github.com/VladimirYalumov/resty/requests"
	"github.com/VladimirYalumov/resty/responses"
	"github.com/VladimirYalumov/tracer"
	"net/http"
)

var additionalMiddlewares []middleware.Middleware

func Init(mm ...middleware.Middleware) {
	additionalMiddlewares = make([]middleware.Middleware, len(mm)+2)
	for i := len(mm) - 1; i != 0; i-- {
		additionalMiddlewares = append(additionalMiddlewares, mm[i])
	}
	additionalMiddlewares = append(additionalMiddlewares, &middleware.RequestValidate{})
}

type handler struct {
	*cors.Cors
	log *logger.Logger

	endpoints map[endpointKey]*endpoint
}

func NewHandler(log *logger.Logger) *handler {
	return &handler{
		log: log,
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

	req := CheckAction(r, e.request(), w)
	if req == nil {
		return
	}

	resp, httpCode := e.Action(ctx, req)
	w.WriteHeader(httpCode)
	if err := resp.PrepareResponse(w); err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		_, _ = w.Write([]byte{})
	}
	return
}

func (h *handler) Endpoint(
	method,
	path string,
	request func() requests.Request,
	action func(ctx context.Context, req requests.Request) (responses.Response, int),
	mm ...string,
) {
	key := endpointKey{path, method}
	h.endpoints[key] = &endpoint{method: method, Action: action, request: request}
	for _, m := range mm {
		h.endpoints[key].middlewares[m] = true
	}
	h.endpoints[key].middlewares[middleware.KeyRequestValidate] = true
	h.endpoints[key].middlewares[middleware.KeyRequestValidate] = true
}
