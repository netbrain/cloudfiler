package web

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
)

type AppError struct {
	status int
	error  error
	stack  []byte
}

func Error(v interface{}) *AppError {
	err := &AppError{
		status: http.StatusInternalServerError,
		stack:  debug.Stack(),
	}

	switch t := v.(type) {
	default:
		err.Textf("%v", t)
	case int:
		err.status = t
	case error:
		err.error = t
	}
	return err
}

func (e *AppError) Text(message string) {
	e.Textf(message, []interface{}{})
}

func (e *AppError) Textf(format string, v ...interface{}) {
	e.error = fmt.Errorf(format, v...)
}

func (e *AppError) Error() string {
	if e.error == nil {
		return strconv.Itoa(e.status)
	}
	return e.error.Error()
}

func (e *AppError) Status() int {
	return e.status
}

func (e *AppError) Stack() []byte {
	return e.stack
}
