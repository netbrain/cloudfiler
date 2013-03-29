package fvgo

import (
	"log"
	"strings"
)

type Form struct {
	attrs  map[string]string
	fields map[string]*field
}

func (f *Form) addField(fi *field) {
	if f.fields == nil {
		f.fields = make(map[string]*field)
	}
	if fi == nil {
		panic("nil input field")
	}

	if len(fi.name()) != 0 {
		f.fields[fi.name()] = fi
	}
}

func (f *Form) Validate() (bool, map[string][]error) {
	log.Printf("Validating form")
	var errs map[string][]error
	for _, field := range f.fields {
		if fieldErrs := field.Validate(); fieldErrs != nil {
			if errs == nil {
				errs = make(map[string][]error)
			}
			errs[field.name()] = fieldErrs
		}
	}
	return len(errs) == 0, errs
}

func (f *Form) AddFormValues(values map[string][]string) {
	log.Println("Adding form values")
	for key, val := range values {
		f.AddFormValue(key, val...)
	}
	log.Println("Done adding form values")
}

func (f *Form) AddFormValue(key string, val ...string) {
	if _, exist := f.fields[key]; exist {
		log.Printf("%s => %s", key, val)
		f.fields[key].value = []string(val)
	}
}

func (f *Form) Action() string {
	return f.attrs["action"]
}

func (f *Form) Method() string {
	if method, exist := f.attrs["method"]; exist {
		return strings.ToUpper(method)
	}
	return "GET" //return default
}

func (f *Form) IsMultipart() bool {
	return f.attrs["enctype"] == "multipart/form-data"
}

func (f *Form) Field(name string) *field {
	return f.fields[name]
}

func (f *Form) Clone() *Form {
	form := (*f).clone()
	fields := form.fields
	form.fields = nil //important
	for _, fi := range fields {
		newField := fi.Clone()
		form.addField(newField)
	}
	return form
}

func (f Form) clone() *Form {
	return &f
}
