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
	middlewares map[string]bool
	action      func(ctx context.Context, req requests.Request) (responses.Response, int)
	request     requests.Request
}
