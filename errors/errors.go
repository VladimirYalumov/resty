package errors

import (
	"encoding/json"
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

type CustomErrors struct {
	HttpCode    int    `json:"httpCode"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

var customErrorsMap map[int32]CustomErrors

func Init(additionalErrorsMap map[int32]CustomErrors) {
	customErrorsMap = make(map[int32]CustomErrors)
	customErrorsMap[ErrorUnableAddData] = CustomErrors{http.StatusInternalServerError, "Internal error", "Unable to add data to the table"}
	customErrorsMap[ErrorUnableSendMessage] = CustomErrors{http.StatusBadRequest, "Send message error", "Unable to send message"}
	customErrorsMap[ErrorInvalidRequest] = CustomErrors{http.StatusNotAcceptable, "Invalid request", "Request does not meet the requirements"}
	customErrorsMap[ErrorCustomError] = CustomErrors{http.StatusBadRequest, "", "Custom api_errors that will continue the flow"}
	customErrorsMap[ErrorIncorrectVerifyCode] = CustomErrors{http.StatusForbidden, "Incorrect verify code", "Incorrect verify code or code is no longer in redis"}
	customErrorsMap[ErrorUnableGetData] = CustomErrors{http.StatusInternalServerError, "Internal error", "Unable to get data from the table"}
	customErrorsMap[ErrorInvalidAccess] = CustomErrors{http.StatusForbidden, "Access denied", "Invalid client or user not a creator"}
	customErrorsMap[ErrorUserNotFound] = CustomErrors{http.StatusNotFound, "User not found", "User not found in db"}
	customErrorsMap[ErrorUserIsNotVerify] = CustomErrors{http.StatusForbidden, "User is not verify", "User in db is not active"}
	customErrorsMap[ErrorUserUnauthorized] = CustomErrors{http.StatusUnauthorized, "User is unauthorized", "User has not an auth token"}
	customErrorsMap[ErrorCritical] = CustomErrors{http.StatusInternalServerError, "critical error", "Some error in code"}
	customErrorsMap[ErrorNotFound] = CustomErrors{http.StatusNotFound, "Not found", "Something not found"}

	for k, v := range additionalErrorsMap {
		customErrorsMap[k] = v
	}
}

func GetCustomError(w http.ResponseWriter, msg string, code int32) {
	errorResponse := new(responses.ErrorResponse)
	if msg == "" {
		errorResponse.Message = customErrorsMap[code].Message
	} else {
		errorResponse.Message = msg
	}
	errorResponse.Code = code
	w.WriteHeader(customErrorsMap[code].HttpCode)
	_ = json.NewEncoder(w).Encode(&errorResponse)
}
