package web

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
)

func TestErrorWithStringType(t *testing.T) {
	text := "test-error"
	err := Error(text)
	if err.Error() != text {
		t.Fatalf("Expected '%s' but got '%s'", text, err.Error())
	}
}

func TestErrorWithErrorType(t *testing.T) {
	text := "test-error"
	e := errors.New(text)
	err := Error(e)
	if err.Error() != text {
		t.Fatalf("Expected '%s' but got '%s'", text, err.Error())
	}
}

func TestErrorWitIntType(t *testing.T) {
	text := strconv.Itoa(http.StatusNotFound)
	err := Error(http.StatusNotFound)
	if err.Error() != text {
		t.Fatalf("Expected '%s' but got '%s'", text, err.Error())
	}
}

func TestErrorStatus(t *testing.T) {
	err := Error(http.StatusNotFound)
	if err.Status() != http.StatusNotFound {
		t.Fatalf("Expected '%v' but got '%v'", http.StatusNotFound, err.Status())
	}
}
