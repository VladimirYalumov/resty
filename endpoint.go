package resty

import (
	"context"
	"github.com/VladimirYalumov/resty/responses"
	"resty/requests"
)

type endpointKey struct {
	path   string
	method string
}

type endpoint struct {
	method  string
	Action  func(ctx context.Context, req requests.Request) (responses.Response, int)
	request requests.Request

	middlewares map[string]bool
}
