package restapi

import (
	"fmt"
)

// ResponseError represents an Error to a REST API call
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Request returned an error. %+v", *e)
}
