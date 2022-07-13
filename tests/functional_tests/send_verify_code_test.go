package functional_tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"resty/config"
	"resty/orm"
	"resty/routing"
	"resty/routing/cusotom_errors"
	"resty/routing/requests"
	"resty/routing/responses"
	"strconv"
	"testing"
)

var clientSendVerifyCode = orm.Client{
	Name: "Deewave",
	Key:  "deewave",
}

var userSendVerifyCodeRequest = orm.User{
	Email:    "vladimir.ylmv@gmail.com",
	Password: "25d55ad283aa400af464c76d713c07ad",
	Name:     "vova",
}

func TestInitSendVerifyCode(t *testing.T) {
	config.Init()
	err := userSendVerifyCodeRequest.Create()
	if err != nil {
		t.Fatal(err)
	}
	err = clientSendVerifyCode.Create()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendVerifyCodeSuccess(t *testing.T) {
	requestBody := requests.SendVerifyCodeRequest{
		Email: userSendVerifyCodeRequest.Email,
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/send_verify_code", bytes.NewBuffer(requestBodyStr))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Client", clientSendVerifyCode.Key)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SendVerifyCode)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.Code)
	expectedResponse := responses.SendVerifyCodeResponse{
		Success: true,
	}
	var actualResponse responses.SendVerifyCodeResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestSendVerifyCodeWithoutRequiredFields(t *testing.T) {
	requestBody := requests.SendVerifyCodeRequest{}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/send_verify_code", bytes.NewBuffer(requestBodyStr))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Client", clientSendVerifyCode.Key)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SendVerifyCode)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusNotAcceptable, response.Code)
	expectedResponse := responses.ErrorResponse{
		Code:    cusotom_errors.ErrorInvalidRequest,
		Message: routing.CustomErrorsMap.Errors[strconv.Itoa(cusotom_errors.ErrorInvalidRequest)].Message,
	}
	var actualResponse responses.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestSendVerifyCodeWithInvalidRequest(t *testing.T) {
	requestBody := requests.SendVerifyCodeRequest{
		Email: "",
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/send_verify_code", bytes.NewBuffer(requestBodyStr))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Client", clientSendVerifyCode.Key)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SendVerifyCode)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusNotAcceptable, response.Code)
	expectedResponse := responses.ErrorResponse{
		Code:    cusotom_errors.ErrorInvalidRequest,
		Message: routing.CustomErrorsMap.Errors[strconv.Itoa(cusotom_errors.ErrorInvalidRequest)].Message,
	}
	var actualResponse responses.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestSendVerifyCodeWithInvalidClient(t *testing.T) {
	requestBody := requests.SendVerifyCodeRequest{
		Email: userSendVerifyCodeRequest.Email,
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/send_verify_code", bytes.NewBuffer(requestBodyStr))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Client", "test")
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SendVerifyCode)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusForbidden, response.Code)
	expectedResponse := responses.ErrorResponse{
		Code:    cusotom_errors.ErrorInvalidClient,
		Message: routing.CustomErrorsMap.Errors[strconv.Itoa(cusotom_errors.ErrorInvalidClient)].Message,
	}
	var actualResponse responses.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestSendVerifyCodeWithNonExistentEmail(t *testing.T) {
	requestBody := requests.SendVerifyCodeRequest{
		Email: "test@mail.ru",
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/send_verify_code", bytes.NewBuffer(requestBodyStr))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Client", clientSendVerifyCode.Key)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SendVerifyCode)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusNotFound, response.Code)
	expectedResponse := responses.ErrorResponse{
		Code:    cusotom_errors.ErrorUserNotFound,
		Message: "User with email test@mail.ru not found",
	}
	var actualResponse responses.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCloseSendVerifyCode(t *testing.T) {
	err := clientSendVerifyCode.Delete()
	if err != nil {
		t.Fatal(err)
	}
	err = userSendVerifyCodeRequest.Delete()
	if err != nil {
		t.Fatal(err)
	}
	config.Close()
}
