package middleware

import (
	"github.com/VladimirYalumov/resty/errors"
	"github.com/VladimirYalumov/resty/requests"
	"net/http"
)

const KeyRequestInit = "request_init"

type RequestInit struct {
	next Middleware
	r    *http.Request
}

func NewRequestInit(r *http.Request) *RequestInit {
	return &RequestInit{r: r}
}

func (r *RequestInit) Execute(req requests.Request) (int32, string) {
	val, ok := req.Middlewares()[r.getKey()]
	if !ok || !val {
		return r.next.Execute(req)
	}
	if err := req.Set(r.r); err != nil {
		return errors.ErrorInvalidRequest, ""
	}

	return r.next.Execute(req)
}

func (r *RequestInit) SetNext(next Middleware) {
	r.next = next
}

func (r *RequestInit) getKey() string {
	return KeyRequestInit
}
