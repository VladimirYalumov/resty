package requests

type GetUserRequest struct {
	Id     int
	Client string
	Token  string
}

func (request GetUserRequest) ValidateRequest() (bool, string) {
	if request.Id == 0 {
		return false, ""
	}
	return true, ""
}

func (request GetUserRequest) GetClient() string {
	return request.Client
}
