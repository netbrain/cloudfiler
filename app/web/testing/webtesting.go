package testing

import (
	. "github.com/netbrain/cloudfiler/app/web"
	"net/http"
	"net/http/httptest"
)

//This file contains shared testing functions

func CreateReqContext(method, path string, parameters map[string][]string) (*Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, path, nil)
	switch method {
	case "GET":
		for key, slice := range parameters {
			for _, val := range slice {
				r.URL.Query().Add(key, val)
			}
		}
	case "POST":
		r.Form = parameters
	default:
		panic("unknown method")
	}

	ctx := NewContext(w, r)
	return ctx, w
}
