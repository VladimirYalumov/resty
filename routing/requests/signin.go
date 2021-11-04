package requests

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Client   string
}

func (request SignInRequest) ValidateRequest() (bool, string) {
	if request.Email == "" || request.Password == "" {
		return false, ""
	}
	return true, ""
}

func (request SignInRequest) GetClient() string {
	return request.Client
}
