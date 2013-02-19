package fvgo

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	requiredErr = "Field '%v' is required"
	emailErr    = "Field '%v' is not an email"
	patternErr  = "Field '%v' doesn't match the pattern '%v'"
)

type field struct {
	tag   string
	attrs map[string]string
	value []string
	regex *regexp.Regexp
}

func NewField(tag string, attrs map[string]string) *field {

	field := &field{
		tag:   tag,
		attrs: attrs,
	}

	if len(field.name()) == 0 && field.ftype() != "submit" {
		panic(fmt.Sprintf("Field '%#v' has zero length name", field))
	}

	if field.hasPattern() {
		field.regex = regexp.MustCompile(field.pattern())
	}

	return field
}

func (f *field) Validate() []error {
	var errs []error

	if f.required() || f.hasValueContent() {
		if f.required() && !f.hasValueContent() {
			errs = append(errs, validateErr(requiredErr, f.name()))
		}

		if f.isEmail() && !f.validateEmail() {
			errs = append(errs, validateErr(emailErr, f.name()))
		}

		if f.hasPattern() {
			ok := f.validatePattern()

			if !ok {
				errs = append(errs, validateErr(patternErr, f.name(), f.pattern()))
			}
		}
	}
	log.Printf("Field '%s' of type '%s' has the following attrs: '%v'", f.name(), f.tag, f.attrs)
	log.Printf("Field '%s' has %v errors: %v where value was: '%v'", f.name(), len(errs), errs, f.value)

	return errs
}

func (f *field) required() bool {
	_, ok := f.attrs["required"]
	return ok
}

func (f *field) ftype() string {
	return f.attrs["type"]
}

func (f *field) isEmail() bool {
	return f.ftype() == "email"
}

func (f *field) validateEmail() bool {
	for _, value := range f.value {
		if !strings.Contains(value, "@") {
			return false
		}
	}
	return true
}

func (f *field) validatePattern() bool {
	for _, value := range f.value {
		if !f.regex.MatchString(value) {
			return false
		}
	}
	return true
}

func (f *field) addValue(val string) {
	f.value = append(f.value, val)
}

func (f *field) setValue(val ...string) {
	f.value = val
}

func (f *field) hasValueContent() bool {
	return len(f.value) > 0
}

func (f *field) isSingleValue() bool {
	return len(f.value) == 1
}

func (f *field) hasPattern() bool {
	return len(f.pattern()) > 0
}

func (f *field) pattern() string {
	p, exists := f.attrs["pattern"]
	if exists {
		if !strings.HasPrefix(p, "^") {
			p = "^" + p
		}
		if !strings.HasSuffix(p, "$") {
			p += "$"
		}
	}
	return p
}

func (f *field) name() string {
	return strings.ToLower(f.attrs["name"])
}

func (f *field) Clone() *field {
	return (*f).clone()
}

func (f field) clone() *field {
	return &f
}

func validateErr(errorf string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(errorf, args...))
}
