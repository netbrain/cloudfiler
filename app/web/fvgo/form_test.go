package fvgo

import (
	"reflect"
	"testing"
)

func TestAddNilField(t *testing.T) {
	f := &Form{}

	defer func() {
		if recover() == nil {
			t.Fatal("Expected panic from adding nil value")
		}
	}()

	f.addField(nil)
	if len(f.fields) > 0 {
		t.Fatal("Should not be able to add nil field")
	}
}

func TestAddField(t *testing.T) {
	form := &Form{}
	attrs := map[string]string{"name": "test"}
	field := NewField("input", attrs)

	form.addField(field)
	if len(form.fields) != 1 {
		t.Fatal("Expected 1 field")
	}
}

func TestMultipartForm(t *testing.T) {
	form := &Form{}

	if form.IsMultipart() {
		t.Fatal("Form should not be multipart.")
	}

	form.attrs = map[string]string{"enctype": "multipart/form-data"}

	if !form.IsMultipart() {
		t.Fatal("Form should be of type multipart")
	}
}

func TestFormClone(t *testing.T) {
	f1 := &Form{
		attrs: map[string]string{
			"action": "#",
			"method": "post",
		},
	}

	f1.addField(NewField("input", map[string]string{
		"type": "text",
		"name": "testfield1",
	}))

	f1.addField(NewField("input", map[string]string{
		"type": "text",
		"name": "testfield2",
	}))

	f2 := f1.Clone()

	fval1 := reflect.ValueOf(f1)
	fval2 := reflect.ValueOf(f2)

	if fval2.IsNil() {
		t.Fatal("Clone is nil")
	}

	if fval1.Pointer() == fval2.Pointer() {
		t.Fatal("Forms have identical pointer")
	}

	if !reflect.DeepEqual(f1, f2) {
		t.Logf("\n%v === %v\n%v === %v\n", f1.attrs, f2.attrs, f1.fields, f2.fields)
		t.Fatal("Forms aren't equal")
	}

	if len(f1.fields) != len(f2.fields) {
		t.Fatal("fields are not equal")
	}

	f1.Field("testfield1").value = []string{"testval"}
	if f2.Field("testfield2").hasValueContent() {
		t.Fatal("f2 should not have value content")
	}
}

func TestValidateEmptyForm(t *testing.T) {
	form := &Form{}
	ok, errs := form.Validate()
	if !ok {
		t.Fatal("should not fail ")
	}

	if errs != nil {
		t.Fatal("errors should be nil map")
	}
}

func TestValidateSimpleForm(t *testing.T) {
	form := &Form{}

	form.addField(NewField("input", map[string]string{
		"name":     "testfield",
		"required": "",
	}))

	form.addField(NewField("input", map[string]string{
		"name":     "testfield2",
		"required": "",
	}))

	ok, errs := form.Validate()
	if ok {
		t.Fatal("should fail validation")
	}

	if len(errs) != 2 {
		t.Fatal("Should have 2 failed validations")
	}
}

func TestFormCloneSuspectingDataLeak(t *testing.T) {
	f1 := &Form{
		attrs: map[string]string{
			"action": "#",
			"method": "post",
		},
	}

	f1.addField(NewField("input", map[string]string{
		"type": "text",
		"name": "testfield1",
	}))

	f2 := f1.Clone()
	field2 := f2.Field("testfield1")
	field2.addValue("some-data-leak")

	f3 := f1.Clone()
	field3 := f3.Field("testfield1")
	if field3.hasValueContent() && field3.value[0] == field2.value[0] {
		t.Fatal("Data leakage!")
	}
}
