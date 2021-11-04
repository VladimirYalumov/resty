package requests

type SendVerifyCodeRequest struct {
	Email  string `json:"email"`
	Client string
}

func (request SendVerifyCodeRequest) ValidateRequest() (bool, string) {
	if request.Email == "" {
		return false, ""
	}
	return true, ""
}

func (request SendVerifyCodeRequest) GetClient() string {
	return request.Client
}
