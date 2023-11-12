package resty

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VladimirYalumov/logger"
	"github.com/VladimirYalumov/tracer"
	"net/http"
	"resty/action"
	"resty/errors"
	"resty/middleware"
	"resty/requests"
	"resty/responses"
	"runtime/debug"
)

var additionalMiddlewares []middleware.Middleware

func Init(mm ...middleware.Middleware) {
	additionalMiddlewares = make([]middleware.Middleware, len(mm)+2)
	for i := len(mm) - 1; i != 0; i-- {
		additionalMiddlewares = append(additionalMiddlewares, mm[i])
	}
	additionalMiddlewares = append(additionalMiddlewares, &middleware.RequestValidate{})
}

type Handler struct {
	*cors.Cors
	log *logger.Logger
}

func NewHandler(log *logger.Logger) *Handler {
	return &Handler{log: log}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func getDeferCatchPanic(log *logger.Logger, w http.ResponseWriter) {
	if rec := recover(); rec != any(nil) {
		logger.Error(
			logger.ToContext(context.Background(), log),
			fmt.Errorf("error: %v", rec), "critical error", "stacktrace", string(debug.Stack()),
		)
		errors.GetCustomError(w, "", errors.ErrorCritical)
		return
	}
}

func CheckAction(r *http.Request, req requests.Request, w http.ResponseWriter) requests.Request {
	currentRequest := &req
	checkRequest := &middleware.RequestCheck{}

	for i, additionalMiddleware := range additionalMiddlewares {
		if i+1 == len(additionalMiddlewares) {
			additionalMiddleware.SetNext(checkRequest)
			break
		}
		additionalMiddleware.SetNext(additionalMiddlewares[i+1])
	}

	initRequest := middleware.NewRequestInit(r)
	initRequest.SetNext(additionalMiddlewares[0])

	code, msg := initRequest.Execute(currentRequest)

	if code != errors.ErrorNoError {
		errors.GetCustomError(w, msg, code)
		return nil
	}

	return *currentRequest
}
