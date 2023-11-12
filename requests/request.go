package requests

import "net/http"

type Request interface {
	Validate() (bool, string)
	Set(r *http.Request) error
	Middlewares() map[string]bool

	GetCustomStruct() any
}
