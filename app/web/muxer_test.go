package web

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestHandler struct {
	closure func(*Context)
}

func (t TestHandler) Action(ctx *Context) interface{} {
	if t.closure != nil {
		t.closure(ctx)
	}
	return "Hello World"
}

var m Muxer
var action = TestHandler.Action
var path = "/some/path"

func initMuxerTest() {
	//reset muxer
	m = NewMuxer()
}

func TestAddHandler(t *testing.T) {
	initMuxerTest()
	m.AddHandler(TestHandler{})

	if len(m.handlers) != 1 {
		t.Fatal("Should be of size 1")
	}

	if _, ok := m.Handler("TestHandler"); !ok {
		t.Fatal("Could not fetch handler")
	}

}

func TestGetHandlerByCaseInsensitiveName(t *testing.T) {
	initMuxerTest()
	m.AddHandler(TestHandler{})

	if _, ok := m.Handler("tEstHanDler"); !ok {
		t.Fatal("Could not fetch handler")
	}

}

func TestAddAnonymousStructHandler(t *testing.T) {
	initMuxerTest()
	myHandler := struct{}{}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected to fail!")
		}
	}()

	m.AddHandler(myHandler)
}

func TestAddAction(t *testing.T) {
	initMuxerTest()
	m.AddAction(path, action)

	if len(m.actions) != 1 {
		t.Fatal("Expected 1 action")
	}

	if _, ok := m.Action("/some/path"); !ok {
		t.Fatal("Expected action returned")
	}
}

func TestHandleActionReturnsData(t *testing.T) {
	initMuxerTest()
	m.AddHandler(TestHandler{})
	m.AddAction(path, action)

	v := reflect.ValueOf(action)
	if result := m.handleAction(v, &Context{}); result != "Hello World" {
		t.Fatalf("Expected 'Hello World' but got: %v", result)
	}
}

func TestHandleActionContextIsInitialized(t *testing.T) {
	initMuxerTest()

	closureCalled := false

	defer func() {
		if !closureCalled {
			t.Fatal("Closure wasn't executed")
		}
	}()

	m.AddHandler(TestHandler{
		closure: func(ctx *Context) {
			closureCalled = true
			if ctx.Request == nil {
				t.Fatal("request not initialized")
			}
		},
	})
	m.AddAction(path, action)

	r, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	m.Handle(w, r)
}

func TestHandleActionWhereNoActionExist(t *testing.T) {
	initMuxerTest()
	m.AddAction(path, action)
	m.Action("/some/unhandled/path")
	invalidAction, _ := m.Action("/some/unhandled/path")
	if result := m.handleAction(invalidAction, &Context{}); result != nil {
		t.Fatalf("Expected 'nil' but got: %v", result)
	}
}

func TestGetActionName(t *testing.T) {
	initMuxerTest()
	name := m.actionName(reflect.ValueOf(TestHandler.Action))
	if name != "action" {
		t.Fatalf("Expected 'action' but got: %s", name)
	}
}

func TestGetHandlerName(t *testing.T) {
	initMuxerTest()
	name := m.handlerName(reflect.ValueOf(TestHandler.Action))
	if name != "test" {
		t.Fatalf("Expected 'test' but got: %s", name)
	}
}

func TestActionPath(t *testing.T) {
	initMuxerTest()
	m.AddHandler(TestHandler{})
	m.AddAction(path, action)
	p, ok := m.ActionPath(action)

	if !ok {
		t.Fatal("Did not find path for input action")
	}

	if p != path {
		t.Fatalf("Expected retrieved ActionPath (%s) to equal %s", p, path)
	}

}
