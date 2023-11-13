package resty

import (
	"context"
	"net/http"
	"resty/requests"
)

type endpointKey struct {
	path   string
	method string
}

type endpoint[T any] struct {
	method  string
	Action  func(ctx context.Context, req requests.Request, w http.ResponseWriter)
	request requests.Request

	middlewares map[string]bool

	data *T
}
