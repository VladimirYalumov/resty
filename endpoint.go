package resty

import (
	"context"
	"github.com/VladimirYalumov/resty/errors"
	"github.com/VladimirYalumov/resty/responses"
	"resty/requests"
)

type endpointKey struct {
	path   string
	method string
}

type endpoint[T any] struct {
	method  string
	Action  func(ctx context.Context, data T, req requests.Request) (responses.Response, errors.CustomError)
	request requests.Request

	middlewares map[string]bool

	data *T
}
