package models

import "encoding/json"

type httpError struct {
	StatusCode int    `json:"status_code"`
	ErrorText  string `json:"error_text"`
}

func NewHttpError(statusCode int, err error) ([]byte, error) {
	httpErr := httpError{
		StatusCode: statusCode,
		ErrorText:  err.Error(),
	}
	payload, err := json.Marshal(httpErr)

	return payload, err
}
