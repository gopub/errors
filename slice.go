package errors

import "strings"

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

func Append(err error, errs ...error) error {
	if s, ok := err.(errorSlice); ok {
		return append(s, errs...)
	}
	s := make(errorSlice, 0, 1+len(errs))
	s = append(s, err)
	s = append(s, errs...)
	return s
}
