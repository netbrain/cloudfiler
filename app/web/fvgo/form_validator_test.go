package fvgo

import (
	"log"
	"net/http"
	"testing"
)

func TestValidateRequestDataSuspectingDataLeak(t *testing.T) {
	fv := NewFormValidator()
	form := &Form{
		attrs: map[string]string{
			"action": "/some/url",
			"method": "post",
		},
	}
	form.addField(NewField("input", map[string]string{
		"type":     "text",
		"name":     "test",
		"required": "",
	}))

	fv.forms = map[string]*Form{
		"/some/url": form,
	}

	req, _ := http.NewRequest("POST", "/some/url", nil)
	req.Form = map[string][]string{
		"test": []string{"testval"},
	}

	if ok, errors := fv.ValidateRequestData(req); !ok {
		t.Fatal(errors)
	}

	req, _ = http.NewRequest("POST", "/some/url", nil)
	req.Form = map[string][]string{}

	if ok, _ := fv.ValidateRequestData(req); ok {
		t.Fatal("Expected errors on second request")
	}

	log.SetPrefix("")

}
