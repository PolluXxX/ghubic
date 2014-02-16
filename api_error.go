package ghubic

import (
	"fmt"
)

type ApiError struct {
	HTTPCode int
	Code     int    `json:"code"`
	Message  string `json:"message"`
	ApiError string `json:"error"`
	ApiDesc  string `json:"error_description"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("HTTP %d : %s (%d)", e.HTTPCode, e.Message, e.Code)
}
