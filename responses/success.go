package responses

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
