package errors

import (
	"github.com/VladimirYalumov/resty/responses"
	"net/http"
)

const ErrorNoError = -1

const ErrorUnableAddData = 0       // ErrorUnableAddData Unable to add data to the table
const ErrorUnableSendMessage = 1   // ErrorUnableSendMessage Unable to send message
const ErrorInvalidRequest = 2      // ErrorInvalidRequest Request does not meet the requirements
const ErrorCustomError = 3         // ErrorCustomError Custom api_errors that will continue the flow
const ErrorIncorrectVerifyCode = 4 // ErrorIncorrectVerifyCode Incorrect verify code or code is no longer in redis
const ErrorUnableGetData = 5       // ErrorUnableGetData Unable to get data from the table
const ErrorInvalidAccess = 6       // ErrorInvalidAccess Invalid access
const ErrorUserNotFound = 7        // ErrorUserNotFound User not found in db
const ErrorUserIsNotVerify = 8     // ErrorUserIsNotVerify User in db is not active
const ErrorUserUnauthorized = 9    // ErrorUserUnauthorized User has not an auth token
const ErrorCritical = 10           // ErrorCritical Some error in code

const ErrorNotFound = 11 // ErrorNotFound Object not found in db

type CustomError struct {
	HttpCode    int    `json:"httpCode"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

var CustomErrorMap map[int32]CustomError

func Init(additionalErrorsMap map[int32]CustomError) {
	CustomErrorMap = make(map[int32]CustomError)
	CustomErrorMap[ErrorUnableAddData] = CustomError{http.StatusInternalServerError, "internal error", "Unable to add data to the table"}
	CustomErrorMap[ErrorUnableSendMessage] = CustomError{http.StatusBadRequest, "send message error", "Unable to send message"}
	CustomErrorMap[ErrorInvalidRequest] = CustomError{http.StatusNotAcceptable, "invalid request", "Request does not meet the requirements"}
	CustomErrorMap[ErrorCustomError] = CustomError{http.StatusBadRequest, "", "custom api_errors that will continue the flow"}
	CustomErrorMap[ErrorIncorrectVerifyCode] = CustomError{http.StatusForbidden, "incorrect verify code", "Incorrect verify code or code is no longer in redis"}
	CustomErrorMap[ErrorUnableGetData] = CustomError{http.StatusInternalServerError, "internal error", "Unable to get data from the table"}
	CustomErrorMap[ErrorInvalidAccess] = CustomError{http.StatusForbidden, "access denied", "Invalid client or user not a creator"}
	CustomErrorMap[ErrorUserNotFound] = CustomError{http.StatusNotFound, "user not found", "User not found in db"}
	CustomErrorMap[ErrorUserIsNotVerify] = CustomError{http.StatusForbidden, "user is not verify", "User in db is not active"}
	CustomErrorMap[ErrorUserUnauthorized] = CustomError{http.StatusUnauthorized, "user is unauthorized", "User has not an auth token"}
	CustomErrorMap[ErrorCritical] = CustomError{http.StatusInternalServerError, "critical error", "Some error in code"}
	CustomErrorMap[ErrorNotFound] = CustomError{http.StatusNotFound, "not found", "Something not found"}

	for k, v := range additionalErrorsMap {
		CustomErrorMap[k] = v
	}
}

func GetCustomError(msg string, code int32) (*responses.ErrorResponse, int) {
	customError, ok := CustomErrorMap[code]
	if !ok {
		return &responses.ErrorResponse{Code: ErrorCustomError, Message: "undefined"}, http.StatusBadRequest
	}

	resp := &responses.ErrorResponse{Code: code}
	if msg == "" {
		resp.Message = customError.Message
	} else {
		resp.Message = msg
	}
	return resp, customError.HttpCode
}
