package errors

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	s := http.StatusText(e.Code)
	if s != "" {
		return s
	}
	return fmt.Sprintf("error: %d", e.Code)
}

func (e *Error) Respond(ctx context.Context, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(e.Code)
	msg := e.Message
	if msg == "" {
		msg = http.StatusText(e.Code)
	}
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Printf("Cannot write: %v", err)
	}
}

func Format(code int, format string, a ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

var rawErrType = reflect.TypeOf(errors.New(""))

func GetCode(err error) int {
	err = Cause(err)
	if reflect.TypeOf(err) == rawErrType {
		return 0
	}

	if err == NotExist || err == sql.ErrNoRows {
		return http.StatusNotFound
	}

	if coder, ok := err.(interface{ Status() int }); ok {
		return coder.Status()
	}

	if coder, ok := err.(interface{ StatusCode() int }); ok {
		return coder.StatusCode()
	}

	if coder, ok := err.(interface{ Code() int }); ok {
		return coder.Code()
	}

	if v := reflect.ValueOf(err); v.Kind() == reflect.Int {
		n := int(v.Int())
		if n > 0 {
			return n
		}
		return 0
	}

	keys := []string{"status", "Status", "status_code", "StatusCode", "statusCode", "code", "Code"}
	i := indirect(err)
	k := reflect.ValueOf(i).Kind()
	if k != reflect.Struct && k != reflect.Map {
		return 0
	}

	b, jErr := json.Marshal(i)
	if jErr != nil {
		log.Printf("Cannot marshal json: %v", err)
		return 0
	}
	var m map[string]interface{}
	jErr = json.Unmarshal(b, &m)
	if jErr != nil {
		log.Printf("Cannot unmarshal json: %v", err)
		return 0
	}

	for _, k := range keys {
		if v, ok := m[k]; ok {
			if i, err := strconv.ParseInt(fmt.Sprint(v), 0, 0); err == nil && i > 0 {
				return int(i)
			}
		}
	}
	return 0
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func IsNotExist(err error) bool {
	return Is(err, NotExist) || Is(err, sql.ErrNoRows) || Is(err, os.ErrNotExist)
}
