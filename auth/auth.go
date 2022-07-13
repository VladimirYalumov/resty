package auth

import (
	"errors"
	"resty/helpers"
	"resty/orm"
	"resty/redis"
	"resty/routing/requests"
)

const ErrorInvalidPassword = "Invalid Password"

func SignUp(request *requests.SignUpRequest) error {
	user := orm.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: helpers.EncryptPassword(request.Password),
	}
	createUserErr := user.Create()
	if createUserErr != nil {
		return createUserErr
	}
	return nil
}

func CheckCode(request requests.VerifyUserRequest) bool {
	code := redis.Client.Get(redis.CreateKey(redis.EmailAuthCode, request.Email)).Val()
	if code == "" {
		return false
	}
	if request.Code != code {
		return false
	}

	return true
}

func VerifyUser(email string) error {
	var user = orm.User{
		Email: email,
	}
	updateUserErr := orm.VerifyUser(&user, true)
	if updateUserErr != nil {
		return updateUserErr
	}

	return nil
}

func SignIn(request *requests.SignInRequest) (bool, error, string) {
	var user orm.User
	find, findUserErr := user.GetByStringValue("email", request.Email)
	if !find {
		return false, findUserErr, ""
	}

	if user.Password != helpers.EncryptPassword(request.Password) {
		errInvalidPassword := errors.New(ErrorInvalidPassword)
		return false, errInvalidPassword, ""
	}

	var client orm.Client
	_, GetClientErr := client.GetByStringValue("key", request.Client)
	if GetClientErr != nil {
		return false, GetClientErr, ""
	}
	token, authTokenError := orm.GetAuthToken(user.Id, client.Id)
	if authTokenError != nil {
		return false, authTokenError, ""
	}
	return true, nil, token.Token
}

func SignOut(userId int, client string) error {
	deleteErr := orm.DeleteToken(userId, client)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func CheckUserByEmail(email string) bool {
	var user orm.User
	find, findUserErr := user.GetByStringValue("email", email)
	if !find && findUserErr == nil {
		return false
	}
	return true
}
