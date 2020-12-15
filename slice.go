package errors

import (
	"errors"
	"strings"
)

type errorSlice []error

var _ error = (errorSlice)(nil)

func (s errorSlice) Error() string {
	var b strings.Builder
	for _, e := range s {
		if b.Len() > 0 {
			b.WriteString("; ")
		}
		b.WriteString(e.Error())
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

func (s errorSlice) Unwrap() error {
	for _, e := range s {
		if u := errors.Unwrap(e); u != nil {
			return u
		}
	}
	return nil
}

func (s errorSlice) Is(err error) bool {
	for _, e := range s {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}

func Append(err error, errs ...error) error {
	if s, ok := err.(errorSlice); ok {
		return append(s, errs...)
	}
	s := make(errorSlice, 0, 1+len(errs))
	s = append(s, err)
	s = append(s, errs...)
	return s
}
