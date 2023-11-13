package resty

import (
	"context"
	"fmt"
	"net/http"
	"resty/errors"
	"resty/middleware"
	"resty/requests"
	"runtime/debug"
)

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
