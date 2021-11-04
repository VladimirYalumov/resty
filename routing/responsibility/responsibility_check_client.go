package responsibility

import (
	"goRestApi_main/orm"
	"goRestApi_main/routing/cusotom_errors"
)

type RequestCheckClient struct {
	next Route
}

func (r *RequestCheckClient) Execute(requestFlow *RequestFlow) cusotom_errors.CurrentError {
	var currentError cusotom_errors.CurrentError
	if !requestFlow.Responsibilities[r.getResponsibility()] {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}
	var client orm.Client
	findClient, getClientErr := client.GetByStringValue("key", requestFlow.Request.GetClient())
	if !findClient {
		if getClientErr == nil {
			currentError.Code = cusotom_errors.ErrorInvalidClient
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

func (r *RequestCheckClient) SetNext(next Route) {
	r.next = next
}

func (r *RequestCheckClient) getResponsibility() string {
	return ResponsibilityCheckClient
}
