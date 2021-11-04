package requests

type SignOutRequest struct {
	UserId int `json:"user_id"`
	Client string
	Token  string
}

func (request SignOutRequest) ValidateRequest() (bool, string) {
	if request.UserId == 0 {
		return false, ""
	}
	return true, ""
}

func (request SignOutRequest) GetClient() string {
	return request.Client
}
