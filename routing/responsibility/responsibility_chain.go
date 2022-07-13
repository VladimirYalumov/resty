package responsibility

import (
	"resty/orm"
	"resty/routing/cusotom_errors"
	"resty/routing/requests"
)

const ResponsibilityCheckBody = "check_body"
const ResponsibilityCheckClient = "check_client"
const ResponsibilityCheckUserByEmail = "check_user_by_email"
const ResponsibilityCheckUserById = "check_user_by_id"
const ResponsibilityCheckVerifyUser = "check_verify_user"
const ResponsibilityCheckAuth = "check_auth"

func GetResponsibilities() map[string]bool {
	return map[string]bool{
		ResponsibilityCheckClient:      true,
		ResponsibilityCheckBody:        true,
		ResponsibilityCheckUserByEmail: false,
		ResponsibilityCheckVerifyUser:  false,
		ResponsibilityCheckUserById:    false,
		ResponsibilityCheckAuth:        false,
	}
}

type Route interface {
	Execute(*RequestFlow) cusotom_errors.CurrentError
	SetNext(Route)
	getResponsibility() string
}

type RequestFlow struct {
	Request          requests.Request
	Responsibilities map[string]bool

	// Add a specific fields so as not to write a method for the interface
	Email  string
	UserId int
	Token  string
	User   orm.User
}

type RequestCheck struct {
	next Route
}

func (r *RequestCheck) Execute(requestFlow *RequestFlow) cusotom_errors.CurrentError {
	return cusotom_errors.CurrentError{Code: cusotom_errors.ErrorNoError, Message: ""}
}

func (r *RequestCheck) SetNext(next Route) {
	r.next = next
}

func (r *RequestCheck) getResponsibility() string {
	return ""
}
