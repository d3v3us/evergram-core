package common

import "fmt"

type ErrorResponse struct {
	Status  string `json:"status"`
	Err     string `json:"err"`
	Message string `json:"message"`
}
type SucceededResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  *T     `json:"result"`
}

func (ex *Extendable[error]) AsErrorResponseWithMsg(message string) *ErrorResponse {
	return &ErrorResponse{
		Status:  "error",
		Err:     fmt.Sprintf("%v", ex.Value),
		Message: message,
	}
}
func (ex *Extendable[error]) AsErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Status:  "error",
		Err:     fmt.Sprintf("%v", ex.Value),
		Message: "",
	}
}
func (ex *Extendable[T]) AsSucceededResponse(message string) *SucceededResponse[T] {
	return &SucceededResponse[T]{
		Status:  "success",
		Message: message,
		Result:  &ex.Value,
	}
}
