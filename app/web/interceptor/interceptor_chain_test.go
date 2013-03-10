package interceptor

import (
	"net/http"
	"testing"
)

type TestInterceptor struct {
	called bool
	retval bool
}

func (t *TestInterceptor) Handle(w http.ResponseWriter, r *http.Request) bool {
	t.called = true
	return t.retval
}

func TestInterceptorStruct(t *testing.T) {
	ti := &TestInterceptor{
		retval: true,
	}

	if !ti.Handle(nil, nil) {
		t.Fatal("did not return true")
	}

	if !ti.called {
		t.Fatal("called property was not set to true")
	}
}

func TestInterceptorChain(t *testing.T) {

	ic := new(InterceptorChain)
	interceptors := []*TestInterceptor{
		&TestInterceptor{
			retval: true,
		},
		&TestInterceptor{
			retval: false,
		},
		&TestInterceptor{
			retval: true,
		},
	}

	for _, i := range interceptors {
		ic.AddInterceptor(i)
	}

	if len(ic.interceptors) != len(interceptors) {
		t.Fatal("added interceptors length does not equal interceptors length")
	}

	ic.ServeHTTP(nil, nil)

	if !interceptors[0].called {
		t.Fatal("Expected this to be called")
	}

	if !interceptors[1].called {
		t.Fatal("Expected this to be called")
	}

	if interceptors[2].called {
		t.Fatal("Expected this to not be called")
	}

}
