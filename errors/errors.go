package errors

import (
	"net/http"
	"resty/responses"
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
	CustomErrorMap[ErrorUnableAddData] = CustomError{http.StatusInternalServerError, "Internal error", "Unable to add data to the table"}
	CustomErrorMap[ErrorUnableSendMessage] = CustomError{http.StatusBadRequest, "Send message error", "Unable to send message"}
	CustomErrorMap[ErrorInvalidRequest] = CustomError{http.StatusNotAcceptable, "Invalid request", "Request does not meet the requirements"}
	CustomErrorMap[ErrorCustomError] = CustomError{http.StatusBadRequest, "", "Custom api_errors that will continue the flow"}
	CustomErrorMap[ErrorIncorrectVerifyCode] = CustomError{http.StatusForbidden, "Incorrect verify code", "Incorrect verify code or code is no longer in redis"}
	CustomErrorMap[ErrorUnableGetData] = CustomError{http.StatusInternalServerError, "Internal error", "Unable to get data from the table"}
	CustomErrorMap[ErrorInvalidAccess] = CustomError{http.StatusForbidden, "Access denied", "Invalid client or user not a creator"}
	CustomErrorMap[ErrorUserNotFound] = CustomError{http.StatusNotFound, "User not found", "User not found in db"}
	CustomErrorMap[ErrorUserIsNotVerify] = CustomError{http.StatusForbidden, "User is not verify", "User in db is not active"}
	CustomErrorMap[ErrorUserUnauthorized] = CustomError{http.StatusUnauthorized, "User is unauthorized", "User has not an auth token"}
	CustomErrorMap[ErrorCritical] = CustomError{http.StatusInternalServerError, "critical error", "Some error in code"}
	CustomErrorMap[ErrorNotFound] = CustomError{http.StatusNotFound, "Not found", "Something not found"}

	for k, v := range additionalErrorsMap {
		CustomErrorMap[k] = v
	}
}

func GetCustomError(w http.ResponseWriter, msg string, code int32) error {
	resp := new(responses.ErrorResponse)
	if msg == "" {
		resp.Message = CustomErrorMap[code].Message
	} else {
		resp.Message = msg
	}
	resp.Code = code
	w.WriteHeader(CustomErrorMap[code].HttpCode)
	return resp.PrepareResponse(w)
}
