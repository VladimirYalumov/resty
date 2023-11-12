package responses

type ErrorResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
