package responsibility

import (
	"goRestApi_main/helpers"
	"goRestApi_main/routing/cusotom_errors"
)

type RequestCheckVerifyUser struct {
	next Route
}

func (r *RequestCheckVerifyUser) Execute(requestFlow *RequestFlow) cusotom_errors.CurrentError {
	var currentError cusotom_errors.CurrentError
	if !requestFlow.Responsibilities[r.getResponsibility()] {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}
	if helpers.IsEmpty(requestFlow.User) {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}
	if !requestFlow.User.Active {
		currentError.Code = cusotom_errors.ErrorUserIsNotVerify
		currentError.Message = ""
		return currentError
	}
	currentError = r.next.Execute(requestFlow)
	return currentError
}

func (r *RequestCheckVerifyUser) SetNext(next Route) {
	r.next = next
}

func (r *RequestCheckVerifyUser) getResponsibility() string {
	return ResponsibilityCheckVerifyUser
}
