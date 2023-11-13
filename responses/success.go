package responses

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (r *SuccessResponse) PrepareResponse(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r)
}
