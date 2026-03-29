package redash

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (err *APIError) Error() string {
	if err == nil {
		return ""
	}
	if err.Message != "" {
		return fmt.Sprintf("redash API error (%d): %s", err.StatusCode, err.Message)
	}
	if err.Body != "" {
		return fmt.Sprintf("redash API error (%d): %s", err.StatusCode, err.Body)
	}
	return fmt.Sprintf("redash API error (%d)", err.StatusCode)
}

func IsStatus(err error, code int) bool {
	typed, ok := err.(*APIError)
	if !ok {
		return false
	}
	return typed.StatusCode == code
}

