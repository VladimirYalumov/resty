package middleware

import (
	"github.com/VladimirYalumov/resty/errors"
	"github.com/VladimirYalumov/resty/requests"
)

type Middleware interface {
	Execute(requests.Request) (int32, string)
	SetNext(Middleware)
	getKey() string
}

type RequestCheck struct {
	next Middleware
}

func (r *RequestCheck) Execute(_ requests.Request) (int32, string) {
	return errors.ErrorNoError, ""
}

func (r *RequestCheck) SetNext(next Middleware) {
	r.next = next
}

func (r *RequestCheck) getKey() string {
	return ""
}
