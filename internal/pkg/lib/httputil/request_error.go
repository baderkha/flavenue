package httputil

import "fmt"

type RequestError struct {
	Err        error
	StatusCode int
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

func NewRequestError(sCode int, err error) *RequestError {
	return &RequestError{
		Err:        err,
		StatusCode: sCode,
	}
}

func NewComposedRequestErr(sCode int, tmpl string) func(args ...any) *RequestError {
	return func(args ...any) *RequestError {
		return NewRequestError(sCode, fmt.Errorf(tmpl, args...))
	}
}
