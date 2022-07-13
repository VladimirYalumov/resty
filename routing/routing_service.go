package routing

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"resty/routing/cusotom_errors"
	"resty/routing/responses"
	"resty/routing/responsibility"
	"strconv"
)

var ErrorResponse responses.ErrorResponse

type CustomErrors struct {
	Errors map[string]struct {
		Code        int    `json:"code"`
		HttpCode    int    `json:"httpCode"`
		Message     string `json:"message"`
		Description string `json:"description"`
	} `json:"errors"`
}

var CustomErrorsMap CustomErrors

func InitErrors() {
	str, err := ioutil.ReadFile("errors.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, &CustomErrorsMap)
	if err != nil {
		panic(err)
	}
}

func getCustomError(currentError cusotom_errors.CurrentError, w http.ResponseWriter) http.ResponseWriter {
	stringCode := strconv.Itoa(currentError.Code)
	if currentError.Message == "" {
		ErrorResponse.Message = CustomErrorsMap.Errors[stringCode].Message
	} else {
		ErrorResponse.Message = currentError.Message
	}
	ErrorResponse.Code = CustomErrorsMap.Errors[stringCode].Code
	w.WriteHeader(CustomErrorsMap.Errors[stringCode].HttpCode)
	json.NewEncoder(w).Encode(&ErrorResponse)
	return w
}

func UnknownMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(405)
	ErrorResponse.Message = "invalid request"
	json.NewEncoder(w).Encode(&ErrorResponse)
}

func CheckAction(w http.ResponseWriter, requestFlow *responsibility.RequestFlow, decodeError error) (bool, http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	checkRequest := &responsibility.RequestCheck{}
	var currentError cusotom_errors.CurrentError
	if decodeError != nil {
		currentError.Code = cusotom_errors.ErrorInvalidRequest
		currentError.Message = ""
		return false, getCustomError(currentError, w)
	}

	checkVerifyUser := &responsibility.RequestCheckVerifyUser{}
	checkVerifyUser.SetNext(checkRequest)

	checkAuth := &responsibility.RequestCheckAuth{}
	checkAuth.SetNext(checkVerifyUser)

	checkUserById := &responsibility.RequestCheckUserById{}
	checkUserById.SetNext(checkAuth)

	checkUserByEmail := &responsibility.RequestCheckUserByEmail{}
	checkUserByEmail.SetNext(checkUserById)

	checkClient := &responsibility.RequestCheckClient{}
	checkClient.SetNext(checkUserByEmail)

	checkBody := &responsibility.RequestCheckBody{}
	checkBody.SetNext(checkClient)

	currentError = checkBody.Execute(requestFlow)

	if currentError.Code != cusotom_errors.ErrorNoError {
		return false, getCustomError(currentError, w)
	}

	return true, w
}
