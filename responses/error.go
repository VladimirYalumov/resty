package responses

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (r *ErrorResponse) PrepareResponse(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r)
}
