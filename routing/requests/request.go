package requests

type Request interface {
	ValidateRequest() (bool, string)
	GetClient() string
}
