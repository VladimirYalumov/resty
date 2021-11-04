package responsibility

import (
	"goRestApi_main/orm"
	"goRestApi_main/routing/cusotom_errors"
)

type RequestCheckUserById struct {
	next Route
}

func (r *RequestCheckUserById) Execute(requestFlow *RequestFlow) cusotom_errors.CurrentError {
	var currentError cusotom_errors.CurrentError
	if !requestFlow.Responsibilities[r.getResponsibility()] {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}

	if requestFlow.UserId == 0 {
		currentError.Code = cusotom_errors.ErrorUserNotFound
		currentError.Message = ""
		return currentError
	}
	var user orm.User
	findUser, getUserErr := user.GetByIntValue("id", requestFlow.UserId)
	if !findUser {
		if getUserErr == nil {
			currentError.Code = cusotom_errors.ErrorUserNotFound
			currentError.Message = ""
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

func (r *RequestCheckUserById) SetNext(next Route) {
	r.next = next
}

func (r *RequestCheckUserById) getResponsibility() string {
	return ResponsibilityCheckUserById
}
