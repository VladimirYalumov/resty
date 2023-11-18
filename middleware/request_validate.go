package middleware

import (
	"fmt"
	"github.com/VladimirYalumov/resty/errors"
	"github.com/VladimirYalumov/resty/requests"
)

const KeyRequestValidate = "request_validate"

type RequestValidate struct {
	next Middleware
}

func (r *RequestValidate) Execute(req requests.Request) (int32, string) {
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

func (r *RequestValidate) GetKey() string {
	return KeyRequestValidate
}
