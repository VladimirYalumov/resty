package resty

import (
	"context"
	"github.com/VladimirYalumov/resty/requests"
	"github.com/VladimirYalumov/resty/responses"
)

type endpointKey struct {
	path   string
	method string
}

type endpoint[R requests.Request] struct {
	method  string
	Action  func(ctx context.Context, req R) (responses.Response, int)
	request R

	middlewares map[string]bool
}
