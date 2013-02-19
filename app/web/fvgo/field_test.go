package fvgo

import (
	"reflect"
	"testing"
)

func TestFieldValidationRequired(t *testing.T) {
	f := NewField("input", map[string]string{
		"name":     "testfield",
		"required": "",
	})

	if err := f.Validate(); err == nil {
		t.Fatal("Should have failed as value is not set")
	}

	f.value = []string{"testval"}

	if f.Validate() != nil {
		t.Fatal("Should not have failed as value is set")
	}
}

func TestFieldValidationEmail(t *testing.T) {
	f := NewField("input", map[string]string{
		"name": "testfield",
		"type": "email",
	})

	if err := f.Validate(); err != nil {
		t.Fatal("Should not have failed as value is not required")
	}

	f.value = []string{"testval"}

	if err := f.Validate(); err == nil {
		t.Fatal("Should have failed as value is not valid email")
	}

	f.value = []string{"test@test.test"}

	if err := f.Validate(); err != nil {
		t.Fatal("Should not have failed as value is valid email")
	}

	f.attrs["required"] = "" //setting field as required
	f.value = []string{""}

	if err := f.Validate(); err == nil {
		t.Fatal("Should have failed as value is not set")
	}

	f.value = []string{"test"}

	if f.Validate() == nil {
		t.Fatal("Should have failed as value is not email")
	}

	f.value = []string{"test@test.test"}

	if f.Validate() != nil {
		t.Fatal("Should not have failed as value is valid email")
	}

}

func TestFieldValidationPattern(t *testing.T) {
	f := NewField("input", map[string]string{
		"name":    "testfield",
		"type":    "text",
		"pattern": "(foo)|(bar)", //allows only foo or bar
	})

	if f.Validate() != nil {
		t.Fatal("Should not have failed as field is not required")
	}

	f.value = []string{"f00"}

	if f.Validate() == nil {
		t.Fatal("Should have failed as value doesn't match pattern")
	}

	f.value = []string{"foo", "bar"}
	if f.Validate() != nil {
		t.Fatal("Should not have failed as value matches pattern")
	}

	f.value = []string{""}
	f.attrs["required"] = ""

	if f.Validate() == nil {
		t.Fatal("Should have failed as field is required")
	}

	f.value = []string{"f00"}

	if f.Validate() == nil {
		t.Fatal("Should have failed as value doesn't match pattern")
	}

	f.value = []string{"foo", "bar"}
	if f.Validate() != nil {
		t.Fatal("Should not have failed as value matches pattern")
	}

}

func TestFieldValidationInvalidPattern(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Expected panic")
		}
	}()
	NewField("input", map[string]string{
		"name":    "testfield",
		"type":    "text",
		"pattern": "(?=^.{8,}$)((?=.*[0-9])|(?=.*[A-Z]]+))(?![.\n])(?=.*[A-Z])(?=.*[a-z]).*$",
	})
}

func TestFieldValidationPatternWithLength(t *testing.T) {
	f := NewField("input", map[string]string{
		"name":    "testfield",
		"type":    "text",
		"pattern": ".{5,10}", //allows only foo or bar
	})

	f.setValue("") //too small

	if f.Validate() == nil {
		t.Fatal("Expected errors")
	}

	f.setValue("1234567") //within range

	if f.Validate() != nil {
		t.Fatal("Expected zero errors")
	}

	f.setValue("12345678910") //within range

	if f.Validate() == nil {
		t.Fatal("Expected errors")
	}
}

func TestFieldInputTypeSubmitDontNeedName(t *testing.T) {
	f := NewField("input", map[string]string{
		"type":  "submit",
		"value": "Submit this!",
	})

	if err := f.Validate(); err != nil {
		t.Fatal("Should not have failed")
	}
}

func TestFieldClone(t *testing.T) {
	f1 := NewField("input", map[string]string{
		"type":  "submit",
		"value": "Submit this!",
	})

	f2 := f1.Clone()

	fval1 := reflect.ValueOf(f1)
	fval2 := reflect.ValueOf(f2)

	if fval2.IsNil() {
		t.Fatal("Clone is nil")
	}

	if fval1.Pointer() == fval2.Pointer() {
		t.Fatal("Fields have identical pointer")
	}

	if !reflect.DeepEqual(f1, f2) {
		t.Fatal("Fields aren't equal")
	}
}

func TestFieldCloneSuspectingDataLeak(t *testing.T) {
	f1 := NewField("input", map[string]string{
		"type":     "text",
		"name":     "data-leak",
		"required": "",
	})

	f2 := f1.Clone()
	f2.addValue("testing-data-leak")

	if errs := f2.Validate(); len(errs) != 0 {
		t.Fatal("Expected zero errors")
	}

	f3 := f1.Clone()

	if f3.hasValueContent() {
		t.Fatal("Should not have value content")
	}

	if errs := f3.Validate(); len(errs) != 1 {
		t.Fatal("Expected one error")
	}

}
