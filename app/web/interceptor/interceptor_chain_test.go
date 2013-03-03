package interceptor

import (
	"net/http"
	"testing"
)

type TestInterceptor struct {
	called bool
	retval bool
}

func (t TestInterceptor) Handle(w http.ResponseWriter, r *http.Request) bool {
	t.called = true
	return t.retval
}

func TestInterceptorChain(t *testing.T) {

	ic := new(InterceptorChain)
	interceptors := []TestInterceptor{
		TestInterceptor{
			retval: true,
		},
		TestInterceptor{
			retval: false,
		},
		TestInterceptor{
			retval: true,
		},
	}

	for _, i := range interceptors {
		ic.AddInterceptor(i)
	}

	ic.ServeHTTP(nil, nil)

	if !interceptors[0].called {
		t.Fatal("Expected this to be called")
	}

	if !interceptors[1].called {
		t.Fatal("Expected this to be called")
	}

	if interceptors[0].called {
		t.Fatal("Expected this to not be called")
	}

}
