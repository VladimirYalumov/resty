package requests

type VerifyUserRequest struct {
	Code   string `json:"code"`
	Email  string `json:"email"`
	Client string
}

func (request VerifyUserRequest) ValidateRequest() (bool, string) {
	if request.Code == "" || request.Email == "" {
		return false, ""
	}
	return true, ""
}

func (request VerifyUserRequest) GetClient() string {
	return request.Client
}
