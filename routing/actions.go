package routing

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"resty/auth"
	"resty/mail"
	"resty/orm"
	"resty/routing/cusotom_errors"
	"resty/routing/requests"
	"resty/routing/responses"
	"resty/routing/responsibility"
	"strconv"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var request requests.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	request.Client = r.Header.Get("Client")
	requestFlow := &responsibility.RequestFlow{
		Request:          request,
		Responsibilities: responsibility.GetResponsibilities(),
	}
	check, w := CheckAction(w, requestFlow, err)
	if !check {
		return
	}

	// Logic
	var currentError cusotom_errors.CurrentError
	if auth.CheckUserByEmail(request.Email) {
		currentError.Code = cusotom_errors.ErrorCustomError
		currentError.Message = fmt.Sprintf("Email %s already exists", request.Email)
		w = getCustomError(currentError, w)
		return
	}
	errInsertToDb := auth.SignUp(&request)
	if errInsertToDb != nil {
		currentError.Code = cusotom_errors.ErrorInvalidRequest
		currentError.Message = ""
		w = getCustomError(currentError, w)
		return
	}
	successSendMail, _ := mail.SendAuthMessage(request.Email)
	if !successSendMail {
		currentError.Code = cusotom_errors.ErrorUnableSendMessage
		currentError.Message = ""
		w = getCustomError(currentError, w)
		return
	}

	w.WriteHeader(200)
	var response = responses.SignUpResponse{Success: true}
	json.NewEncoder(w).Encode(&response)
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	var request requests.VerifyUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	request.Client = r.Header.Get("Client")
	requestFlow := &responsibility.RequestFlow{
		Request:          request,
		Responsibilities: responsibility.GetResponsibilities(),
		Email:            request.Email,
	}
	requestFlow.Responsibilities[responsibility.ResponsibilityCheckUserByEmail] = true
	check, w := CheckAction(w, requestFlow, err)
	if !check {
		return
	}

	// Logic
	var currentError cusotom_errors.CurrentError
	if auth.CheckCode(request) {
		if auth.VerifyUser(request.Email) != nil {
			currentError.Code = cusotom_errors.ErrorUnableAddData
			currentError.Message = ""
			w = getCustomError(currentError, w)
			return
		}
	} else {
		currentError.Code = cusotom_errors.ErrorIncorrectVerifyCode
		currentError.Message = ""
		w = getCustomError(currentError, w)
		return
	}

	w.WriteHeader(200)
	var response = responses.VerifyUserResponse{Success: true}
	json.NewEncoder(w).Encode(&response)
}

func SendVerifyCode(w http.ResponseWriter, r *http.Request) {
	var request requests.SendVerifyCodeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	request.Client = r.Header.Get("Client")
	requestFlow := &responsibility.RequestFlow{
		Request:          request,
		Responsibilities: responsibility.GetResponsibilities(),
		Email:            request.Email,
	}
	requestFlow.Responsibilities[responsibility.ResponsibilityCheckUserByEmail] = true
	check, w := CheckAction(w, requestFlow, err)
	if !check {
		return
	}

	// Logic
	var currentError cusotom_errors.CurrentError
	successSendMail, errSendMail := mail.SendAuthMessage(request.Email)
	if !successSendMail && errSendMail == nil {
		currentError.Code = cusotom_errors.ErrorCustomError
		currentError.Message = fmt.Sprintf("For email %s, the message sending limit has been exceeded. Wait %s minutes", request.Email, mail.CodeLifeTime)
		w = getCustomError(currentError, w)
		return
	}
	if !successSendMail {
		currentError.Code = cusotom_errors.ErrorUnableSendMessage
		currentError.Message = ""
		w = getCustomError(currentError, w)
		return
	}
	w.WriteHeader(200)
	var response = responses.SendVerifyCodeResponse{Success: true}
	json.NewEncoder(w).Encode(&response)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var request requests.SignInRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	request.Client = r.Header.Get("Client")
	requestFlow := &responsibility.RequestFlow{
		Request:          request,
		Responsibilities: responsibility.GetResponsibilities(),
		Email:            request.Email,
	}
	requestFlow.Responsibilities[responsibility.ResponsibilityCheckUserByEmail] = true
	requestFlow.Responsibilities[responsibility.ResponsibilityCheckVerifyUser] = true
	check, w := CheckAction(w, requestFlow, err)
	if !check {
		return
	}

	// Logic
	var currentError cusotom_errors.CurrentError

	successGetToken, getTokenError, token := auth.SignIn(&request)

	if !successGetToken && getTokenError != nil {
		if getTokenError.Error() == auth.ErrorInvalidPassword {
			currentError.Code = cusotom_errors.ErrorCustomError
			currentError.Message = auth.ErrorInvalidPassword
			w = getCustomError(currentError, w)
			return
		}
		currentError.Code = cusotom_errors.ErrorUnableGetData
		currentError.Message = ""
		w = getCustomError(currentError, w)
		return
	}

	w.WriteHeader(200)
	var response = responses.SignInResponse{Success: true, Token: token, UserId: requestFlow.User.Id}
	json.NewEncoder(w).Encode(&response)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	var request requests.SignOutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	request.Token = r.Header.Get("Authorization")
	request.Client = r.Header.Get("Client")
	requestFlow := &responsibility.RequestFlow{
		Request:          request,
		Responsibilities: responsibility.GetResponsibilities(),
		UserId:           request.UserId,
		Token:            request.Token,
	}

	requestFlow.Responsibilities[responsibility.ResponsibilityCheckUserById] = true
	requestFlow.Responsibilities[responsibility.ResponsibilityCheckVerifyUser] = true
	requestFlow.Responsibilities[responsibility.ResponsibilityCheckAuth] = true
	check, w := CheckAction(w, requestFlow, err)
	if !check {
		return
	}

	// Logic
	var currentError cusotom_errors.CurrentError
	deleteToken := auth.SignOut(request.UserId, request.Client)
	if deleteToken != nil {
		currentError.Code = cusotom_errors.ErrorUnableGetData
		currentError.Message = ""
		w = getCustomError(currentError, w)
		return
	}

	w.WriteHeader(200)
	var response = responses.SignOutResponse{Success: true}
	json.NewEncoder(w).Encode(&response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var request requests.GetUserRequest
	var err error
	params := mux.Vars(r)
	request.Token = r.Header.Get("Authorization")
	request.Client = r.Header.Get("Client")
	request.Id, err = strconv.Atoi(params["id"])
	requestFlow := &responsibility.RequestFlow{
		Request:          request,
		Responsibilities: responsibility.GetResponsibilities(),
		UserId:           request.Id,
		Token:            request.Token,
	}

	requestFlow.Responsibilities[responsibility.ResponsibilityCheckUserById] = true
	requestFlow.Responsibilities[responsibility.ResponsibilityCheckVerifyUser] = true
	requestFlow.Responsibilities[responsibility.ResponsibilityCheckAuth] = true
	check, w := CheckAction(w, requestFlow, err)
	if !check {
		return
	}

	// Logic
	var currentError cusotom_errors.CurrentError
	var user orm.User
	_, getUserError := user.GetByIntValue("id", request.Id)
	if getUserError != nil {
		currentError.Code = cusotom_errors.ErrorUnableGetData
		currentError.Message = ""
		w = getCustomError(currentError, w)
		return
	}
	w.WriteHeader(200)
	var response = responses.GetUserResponse{
		UserId:   user.Id,
		UserName: user.Name,
	}
	json.NewEncoder(w).Encode(&response)
}
