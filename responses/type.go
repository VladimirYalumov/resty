package responses

import "net/http"

type Response interface {
	PrepareResponse(w http.ResponseWriter) error
}
