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

type endpoint struct {
	method  string
	Action  func(ctx context.Context, req requests.Request) (responses.Response, int)
	request func() requests.Request

	middlewares map[string]bool
}
