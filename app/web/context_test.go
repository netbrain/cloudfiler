package web

import (
	"errors"
	"reflect"
	"testing"
)

func TestAddFieldValidationError(t *testing.T) {
	ctx := &Context{
		ValidationErrors: make(map[string][]error),
	}
	err := errors.New("test-error")
	ctx.AddFieldValidationError("testfield", err)

	result := ctx.ValidationErrors["testfield"]
	expected := []error{err}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %#v to be equal to %#v", result, expected)
	}
}

func TestAddValidationError(t *testing.T) {
	ctx := &Context{
		ValidationErrors: make(map[string][]error),
	}
	err := errors.New("test-generic-error")
	ctx.AddValidationError(err)

	result := ctx.ValidationErrors[""]
	expected := []error{err}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %#v to be equal to %#v", result, expected)
	}
}

func TestHasValidationErrors(t *testing.T) {
	ctx := &Context{
		ValidationErrors: make(map[string][]error),
	}
	if ctx.HasValidationErrors() {
		t.Fatal("Should not have any errors")
	}
	ctx.AddValidationError(errors.New(""))

	if !ctx.HasValidationErrors() {
		t.Fatal("Should have one error")
	}
}
