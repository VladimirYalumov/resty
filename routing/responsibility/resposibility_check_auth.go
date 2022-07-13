package responsibility

import (
	"resty/helpers"
	"resty/orm"
	"resty/routing/cusotom_errors"
)

type RequestCheckAuth struct {
	next Route
}

func (r *RequestCheckAuth) Execute(requestFlow *RequestFlow) cusotom_errors.CurrentError {
	var currentError cusotom_errors.CurrentError
	if !requestFlow.Responsibilities[r.getResponsibility()] {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}

	if helpers.IsEmpty(requestFlow.User) {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}

	success, err := orm.CheckAuth(requestFlow.UserId, requestFlow.Token, requestFlow.Request.GetClient())

	if !success {
		if err == nil {
			currentError.Code = cusotom_errors.ErrorUserUnauthorized
			currentError.Message = ""
			return currentError
		}
		currentError.Code = cusotom_errors.ErrorUnableGetData
		currentError.Message = ""
		return currentError
	}

	currentError = r.next.Execute(requestFlow)
	return currentError
}

func (r *RequestCheckAuth) SetNext(next Route) {
	r.next = next
}

func (r *RequestCheckAuth) getResponsibility() string {
	return ResponsibilityCheckAuth
}
