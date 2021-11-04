package cusotom_errors

const ErrorNoError = -1
const ErrorUnableAddData = 0       // ErrorUnableAddData Unable to add data to the table
const ErrorUnableSendMessage = 1   // ErrorUnableSendMessage Unable to send message
const ErrorInvalidRequest = 2      // ErrorInvalidRequest Request does not meet the requirements
const ErrorCustomError = 3         // ErrorCustomError Custom errors that will continue the flow
const ErrorIncorrectVerifyCode = 4 // ErrorIncorrectVerifyCode Incorrect verify code or code is no longer in redis
const ErrorUnableGetData = 5       // ErrorUnableGetData Unable to get data from the table
const ErrorInvalidClient = 6       // ErrorInvalidClient Invalid client
const ErrorUserNotFound = 7        // ErrorUserNotFound User not found in db
const ErrorUserIsNotVerify = 8     // ErrorUserIsNotVerify User in db is not active
const ErrorUserUnauthorized = 9    // ErrorUserUnauthorized User has not an auth token

type CurrentError struct {
	Code    int
	Message string
}
