package responses

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
