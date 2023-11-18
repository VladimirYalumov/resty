package resty

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VladimirYalumov/logger"
	"github.com/VladimirYalumov/resty/errors"
	"github.com/VladimirYalumov/resty/middleware"
	"github.com/VladimirYalumov/resty/requests"
	"net/http"
	"runtime/debug"
)

func getDeferCatchPanic(log *logger.Logger, w http.ResponseWriter) {
	if rec := recover(); rec != any(nil) {
		logger.Error(
			logger.ToContext(context.Background(), log),
			fmt.Errorf("error: %v", rec), "critical error", "stacktrace", string(debug.Stack()),
		)
		resp, httpCode := errors.GetCustomError("", errors.ErrorCritical)
		w.WriteHeader(httpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
}

func CheckAction(r *http.Request, req requests.Request, w http.ResponseWriter) requests.Request {
	checkRequest := &middleware.RequestCheck{}
	fmt.Println(additionalMiddlewares)
	for i := 0; i < len(additionalMiddlewares); i++ {
		if i+1 == len(additionalMiddlewares) {
			additionalMiddlewares[i].SetNext(checkRequest)
			break
		}
		additionalMiddlewares[i].SetNext(additionalMiddlewares[i+1])
	}

	initRequest := middleware.NewRequestInit(r)
	initRequest.SetNext(additionalMiddlewares[0])

	code, msg := initRequest.Execute(req)

	if code != errors.ErrorNoError {
		resp, httpCode := errors.GetCustomError(msg, code)
		w.WriteHeader(httpCode)
		_ = json.NewEncoder(w).Encode(resp)
		return nil
	}

	return req
}
