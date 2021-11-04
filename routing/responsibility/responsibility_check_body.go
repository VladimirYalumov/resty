package responsibility

import (
	"fmt"
	"goRestApi_main/routing/cusotom_errors"
)

type RequestCheckBody struct {
	next Route
}

func (r *RequestCheckBody) Execute(requestFlow *RequestFlow) cusotom_errors.CurrentError {
	var currentError cusotom_errors.CurrentError
	if !requestFlow.Responsibilities[r.getResponsibility()] {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}
	valid, field := requestFlow.Request.ValidateRequest()
	if valid {
		currentError = r.next.Execute(requestFlow)
		return currentError
	}
	if field == "" {
		currentError.Code = cusotom_errors.ErrorInvalidRequest
		currentError.Message = ""
		return currentError
	} else {
		currentError.Code = cusotom_errors.ErrorCustomError
		currentError.Message = fmt.Sprintf("Field %s is missing", field)
		return currentError
	}
}

func (r *RequestCheckBody) SetNext(next Route) {
	r.next = next
}

func (r *RequestCheckBody) getResponsibility() string {
	return ResponsibilityCheckBody
}
