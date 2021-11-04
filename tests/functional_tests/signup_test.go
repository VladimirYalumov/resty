package functional_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"goRestApi_main/config"
	"goRestApi_main/helpers"
	"goRestApi_main/orm"
	"goRestApi_main/routing"
	"goRestApi_main/routing/cusotom_errors"
	"goRestApi_main/routing/requests"
	"goRestApi_main/routing/responses"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var user1 = orm.User{
	Email: "vladimir.ylmv@gmail.com",
}

var clientSignUp = orm.Client{
	Name: "Deewave",
	Key:  "deewave",
}

func TestInitSignUp(t *testing.T) {
	config.Init()
	err := clientSignUp.Create()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignUpSuccess(t *testing.T) {
	requestBody := requests.SignUpRequest{
		Name:     "test",
		Email:    user1.Email,
		Password: "12345678",
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("PUT", "/signup", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", clientSignUp.Key)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SignUp)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.Code)
	expectedResponse := responses.SignUpResponse{
		Success: true,
	}
	var actualResponse responses.SignUpResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
	err = user1.Delete()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignUpWithoutRequiredFields(t *testing.T) {
	requestBody := requests.SignUpRequest{
		Name: "test",
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("PUT", "/signup", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", clientSignUp.Key)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SignUp)
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

func TestSignUpWithInvalidRequest(t *testing.T) {
	requestBody := requests.SignUpRequest{
		Name:     "test",
		Email:    "",
		Password: "12345678",
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("PUT", "/signup", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", clientSignUp.Key)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SignUp)
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

func TestSignUpWithInvalidclientSignUp(t *testing.T) {
	requestBody := requests.SignUpRequest{
		Name:     "test",
		Email:    user1.Email,
		Password: "12345678",
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("PUT", "/signup", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", "test")
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SignUp)
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

func TestSignUpTryToCreatedDuplicateUser(t *testing.T) {
	var user2 = orm.User{
		Name:     "test",
		Email:    user1.Email,
		Password: helpers.EncryptPassword("12345678"),
	}
	errCreateUser := user2.Create()
	if errCreateUser != nil {
		t.Fatal(errCreateUser)
	}
	requestBody := requests.SignUpRequest{
		Name:     "test",
		Email:    user1.Email,
		Password: "12345678",
	}
	requestBodyStr, _ := json.Marshal(requestBody)
	request, err := http.NewRequest("PUT", "/signup", bytes.NewBuffer(requestBodyStr))
	request.Header.Set("Client", clientSignUp.Key)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(routing.SignUp)
	handler.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, response.Code)
	expectedResponse := responses.ErrorResponse{
		Code:    cusotom_errors.ErrorCustomError,
		Message: fmt.Sprintf("Email %s already exists", user2.Email),
	}
	var actualResponse responses.ErrorResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedResponse, actualResponse)
	err = user2.Delete()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCloseSignUp(t *testing.T) {
	err := clientSignUp.Delete()
	if err != nil {
		t.Fatal(err)
	}
	err = user1.Delete()
	if err != nil {
		t.Fatal(err)
	}
	config.Close()
}
