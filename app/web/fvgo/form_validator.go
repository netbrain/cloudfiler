package fvgo

import (
	"log"
	"net/http"
)

type FormValidator struct {
	forms map[string]*Form
}

func NewFormValidator() *FormValidator {
	return &FormValidator{
		forms: make(map[string]*Form),
	}
}

func (v *FormValidator) ValidateRequestData(r *http.Request) (bool, map[string][]error) {
	var validFormData = true
	var errs map[string][]error

	log.Printf("Validating against url path: %s, with Method: %s", r.URL.Path, r.Method)
	if form, present := v.forms[r.URL.Path]; present &&
		form.Method() == r.Method {
		log.Println("Found form to validate against")

		form = form.Clone()
		if form.IsMultipart() {
			log.Println("Form is multipart")
			r.ParseMultipartForm(10 << 20) //10mb
			form.AddFormValues(r.MultipartForm.Value)
			for key, val := range r.MultipartForm.File {
				var filenames []string
				for _, fh := range val {
					filenames = append(filenames, fh.Filename)
				}
				form.AddFormValue(key, filenames...)
			}
		} else {
			log.Println("Form is not multipart")
			r.ParseForm()
			form.AddFormValues(r.Form)
		}
		validFormData, errs = form.Validate()
		if validFormData {
			log.Println("Form data passes validation")
		} else {
			log.Println("Form data does not pass validation")
		}
	}

	return validFormData, errs
}

func (v *FormValidator) ParseTemplate(path string) {
	//get form validation data
	parsedForms, err := Parse(path)
	if err != nil {
		panic(err)
	}

	if len(parsedForms) > 0 {
		log.Printf("Found %d form definition(s)", len(parsedForms))
	}

	//TODO form will collide like this
	for _, form := range parsedForms {
		log.Printf("Form has %v fields", len(form.fields))
		action := form.Action()
		if _, present := v.forms[action]; present {
			panic("Form already registered for this action path")
		}
		v.forms[action] = form
	}
}
