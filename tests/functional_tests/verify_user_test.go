package functional_tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"goRestApi_main/config"
	"goRestApi_main/orm"
	"goRestApi_main/redis"
	"goRestApi_main/routing"
	"goRestApi_main/routing/cusotom_errors"
	"goRestApi_main/routing/requests"
	"goRestApi_main/routing/responses"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

const code = "1234"

var clientVerifyUser = orm.Client{
	Name: "Deewave",
	Key:  "deewave",
}

var userVerifyUserRequest = orm.User{
	Email:    "vladimir.ylmv@gmail.com",
	Password: "25d55ad283aa400af464c76d713c07ad",
	Name:     "vova",
}

func TestInitVerifyUser(t *testing.T) {
	config.Init()
	err := userVerifyUserRequest.Create()
	if err != nil {
		t.Fatal(err)
	}
	timeLimit, timeError := time.ParseDuration("10m")
	if timeError != nil {
		t.Fatal(err)
	}
	redis.RedisClient.Set(redis.CreateKey(redis.REDIS_EMAIL_AUTH_CODE, userVerifyUserRequest.Email), code, timeLimit)
	err = clientVerifyUser.Create()
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifyUserSuccess(t *testing.T) {
	requestBody := requests.VerifyUserRequest{
		Code:  code,
		Email: userVerifyUserRequest.Email,
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/verify_user", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", clientVerifyUser.Key)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.VerifyUser)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.Code)
	expectedResponse := responses.VerifyUserResponse{
		Success: true,
	}
	var actualResponse responses.VerifyUserResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestVerifyUserWithoutRequiredFields(t *testing.T) {
	requestBody := requests.VerifyUserRequest{}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/verify_user", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", clientVerifyUser.Key)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.VerifyUser)
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

func TestVerifyUserWithInvalidRequest(t *testing.T) {
	requestBody := requests.VerifyUserRequest{
		Code:  code,
		Email: "",
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/verify_user", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", clientVerifyUser.Key)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.VerifyUser)
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

func TestVerifyUserWithInvalidClient(t *testing.T) {
	requestBody := requests.VerifyUserRequest{
		Code:  code,
		Email: userVerifyUserRequest.Email,
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/verify_user", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", "test")
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.VerifyUser)
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

func TestVerifyUserWithInvalidCode(t *testing.T) {
	requestBody := requests.VerifyUserRequest{
		Code:  "3333",
		Email: userVerifyUserRequest.Email,
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("POST", "/verify_user", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", clientVerifyUser.Key)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.VerifyUser)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusForbidden, response.Code)
	expectedResponse := responses.ErrorResponse{
		Code:    cusotom_errors.ErrorIncorrectVerifyCode,
		Message: routing.CustomErrorsMap.Errors[strconv.Itoa(cusotom_errors.ErrorIncorrectVerifyCode)].Message,
	}
	var actualResponse responses.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestCloseVerifyUser(t *testing.T) {
	err := clientVerifyUser.Delete()
	if err != nil {
		t.Fatal(err)
	}
	err = userVerifyUserRequest.Delete()
	if err != nil {
		t.Fatal(err)
	}
	redis.RedisClient.Del(redis.CreateKey(redis.REDIS_EMAIL_AUTH_CODE, userVerifyUserRequest.Email))
	config.Close()
}
