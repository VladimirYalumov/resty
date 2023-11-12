package middleware

import (
	"fmt"
	"resty/errors"
	"resty/requests"
)

const KeyRequestValidate = "request_validate"

type RequestValidate struct {
	next Middleware
}

func (r *RequestValidate) Execute(req *requests.Request) (int32, string) {
	val, ok := req.Middlewares()[r.getKey()]
	if !ok || !val {
		return r.next.Execute(req)
	}
	valid, field := req.Validate()
	if valid {
		return r.next.Execute(req)
	}
	if field == "" {
		return errors.ErrorInvalidRequest, ""
	} else {
		return errors.ErrorCustomError, fmt.Sprintf("field %s is missing", field)
	}
}

func (r *RequestValidate) SetNext(next Middleware) {
	r.next = next
}

func (r *RequestValidate) getKey() string {
	return KeyRequestValidate
}
