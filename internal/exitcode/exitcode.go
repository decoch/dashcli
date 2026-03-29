package exitcode

import (
	"errors"
	"fmt"
)

const (
	CodeSuccess = 0
	CodeRuntime = 1
	CodeUsage   = 2
)

type Error struct {
	Code int
	Err  error
}

func (err *Error) Error() string {
	if err.Err == nil {
		return "unknown error"
	}
	return err.Err.Error()
}

func (err *Error) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

func Usagef(format string, args ...any) error {
	return &Error{Code: CodeUsage, Err: fmt.Errorf(format, args...)}
}

func Runtimef(format string, args ...any) error {
	return &Error{Code: CodeRuntime, Err: fmt.Errorf(format, args...)}
}

func WrapUsage(err error) error {
	if err == nil {
		return nil
	}
	return &Error{Code: CodeUsage, Err: err}
}

func WrapRuntime(err error) error {
	if err == nil {
		return nil
	}
	return &Error{Code: CodeRuntime, Err: err}
}

func Code(err error) int {
	if err == nil {
		return CodeSuccess
	}
	var typed *Error
	if errors.As(err, &typed) {
		if typed.Code != 0 {
			return typed.Code
		}
	}
	return CodeRuntime
}
