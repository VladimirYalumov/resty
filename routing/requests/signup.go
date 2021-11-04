package requests

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Client   string
}

func (request SignUpRequest) ValidateRequest() (bool, string) {
	if request.Name == "" || request.Email == "" || request.Password == "" {
		return false, ""
	}
	return true, ""
}

func (request SignUpRequest) GetClient() string {
	return request.Client
}
