package responses

type SignInResponse struct {
	Success bool   `json:"success"`
	UserId  int    `json:"user"`
	Token   string `json:"token"`
}
