package responsibility

import (
	"fmt"
	"goRestApi_main/orm"
	"goRestApi_main/routing/cusotom_errors"
)

type RequestCheckUserByEmail struct {
	next Route
}

func (r *RequestCheckUserByEmail) Execute(requestFlow *RequestFlow) cusotom_errors.CurrentError {
	var currentError cusotom_errors.CurrentError
	if !requestFlow.Responsibilities[r.getResponsibility()] {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}

	if requestFlow.Email == "" {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}
	var user orm.User
	findUser, getUserErr := user.GetByStringValue("email", requestFlow.Email)
	if !findUser {
		if getUserErr == nil {
			currentError.Code = cusotom_errors.ErrorUserNotFound
			currentError.Message = fmt.Sprintf("User with email %s not found", requestFlow.Email)
			return currentError
		}
		currentError.Code = cusotom_errors.ErrorUnableGetData
		currentError.Message = ""
		return currentError
	}
	requestFlow.User = user
	currentError = r.next.Execute(requestFlow)
	return currentError
}

func (r *RequestCheckUserByEmail) SetNext(next Route) {
	r.next = next
}

func (r *RequestCheckUserByEmail) getResponsibility() string {
	return ResponsibilityCheckUserByEmail
}
