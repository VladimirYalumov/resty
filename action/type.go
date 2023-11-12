package action

import (
	"context"
	"net/http"
	"resty/middleware"
	"resty/requests"
)

var endpoints map[endpointKey]*endpoint

type endpointKey struct {
	path   string
	method string
}

type endpoint struct {
	method  string
	Action  func(ctx context.Context, req requests.Request, w http.ResponseWriter)
	request requests.Request

	middlewares map[string]bool
}

func New(method, path string, request requests.Request, action func(ctx context.Context, req requests.Request, w http.ResponseWriter), mm ...string) {
	key := endpointKey{path, method}
	endpoints[key] = &endpoint{method: method, Action: action, request: request}
	for _, m := range mm {
		endpoints[key].middlewares[m] = true
	}
	endpoints[key].middlewares[middleware.KeyRequestValidate] = true
	endpoints[key].middlewares[middleware.KeyRequestValidate] = true
}

func GetEndpoint(method, path string) *endpoint {
	e, ok := endpoints[endpointKey{path: path, method: method}]
	if !ok {
		return nil
	}
	return e
}

func (e *endpoint) Request() requests.Request {
	return e.request
}
