package responses

type GetUserResponse struct {
	UserId       int      `json:"id"`
	UserName     string   `json:"name"`
	UserFeatures []string `json:"features"`
	UserImage    string   `json:"image"`
}
