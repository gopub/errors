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
	a, ok := err.(errorSlice)
	if ok {
		return append(a, errs...)
	}
	a = make(errorSlice, 0, 1+len(errs))
	if err != nil {
		a = append(a, err)
	}
	for _, er := range errs {
		if er != nil {
			a = append(a, errs...)
		}
	}
	if len(a) == 0 {
		return nil
	}
	return a
}
