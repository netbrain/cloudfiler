package web

import (
	"bufio"
	"fmt"
	. "github.com/netbrain/cloudfiler/app/web/auth"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	. "strings"
)

var ctxType = reflect.TypeOf(&Context{})

type Muxer struct {
	authenticator Authenticator
	handlers      map[string]reflect.Value
	actions       map[string]reflect.Value
	actionPaths   map[reflect.Value]string
}

func NewMuxer(authenticator Authenticator) Muxer {
	return Muxer{
		authenticator: authenticator,
		handlers:      make(map[string]reflect.Value),
		actions:       make(map[string]reflect.Value),
		actionPaths:   make(map[reflect.Value]string),
	}
}

func (m Muxer) AddHandler(c interface{}) {
	t := reflect.TypeOf(c)
	if t.Kind() != reflect.Struct {
		panic("Expected struct")
	}
	handlerName := ToLower(t.Name())
	if len(handlerName) == 0 {
		panic("Anonymous struct not allowed")
	}
	m.handlers[handlerName] = reflect.ValueOf(c)
}

func (m Muxer) Handler(name string) (v reflect.Value, b bool) {
	v, b = m.handlers[ToLower(name)]
	return
}

func (m Muxer) AddAction(path string, action interface{}) {
	validateAction(action)
	aValue := reflect.ValueOf(action)
	m.actions[path] = aValue
	m.actionPaths[aValue] = path
}

func (m Muxer) Action(path string) (v reflect.Value, b bool) {
	v, b = m.actions[path]
	return
}

func (m Muxer) ActionPath(action interface{}) (path string, b bool) {
	validateAction(action)
	aValue := reflect.ValueOf(action)
	path, b = m.actionPaths[aValue]
	return
}

func validateAction(action interface{}) {
	t := reflect.TypeOf(action)
	if t.Kind() != reflect.Func {
		panic("Expected func")
	}

	if t.NumIn() != 2 {
		panic("Expected 2 input parameters")
	}

	if t.NumOut() != 1 {
		panic("Expected 1 output parameter")
	}

	if t.Out(0).Kind() != reflect.Interface {
		panic("Expected interface{} as output parameter")
	}

	if !t.In(1).AssignableTo(ctxType) {
		panic(fmt.Sprintf("Expected %v as input argument, not: %v", ctxType, t.In(1)))
	}
}

func (m Muxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bw := bufio.NewWriter(w)
	//init failure handling
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, fmt.Sprintf("PANIC: %s - %s", r, debug.Stack()), http.StatusInternalServerError)
			log.Printf("PANIC: %s - %s", r, debug.Stack())
		} else {
			bw.Flush()
		}
	}()
	path := r.URL.Path
	action, ok := m.Action(path)

	if ok {
		//authenticate request
		if m.authenticator.IsAuthorized(r) || path == "/auth/login" {

			//get view for action
			view := m.getViewForAction(action)

			ctx := NewContext(w, r)

			//validate any input form
			if !ctx.HasValidationErrors() {
				//handle action
				ctx.Data = m.handleAction(action, ctx)
			}

			switch d := ctx.Data.(type) {
			case *AppError:
				m.handleAppError(d, w, bw)
			default:
				if ctx.IsRedirected() {
					m.handleRedirect(ctx)
				} else {
					//handle view
					RenderView(view, bw, ctx)
				}
			}
		} else {
			//redirect to login page
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
		}
	} else {
		log.Printf("No handler for path: %s returning status: %v", path, http.StatusNotFound)
		view := "error/404"
		w.WriteHeader(http.StatusNotFound)
		RenderView(view, w, nil)
	}
}

func (m Muxer) handleAction(action reflect.Value, ctx *Context) interface{} {
	if action.IsValid() {
		handlerName := ToLower(action.Type().In(0).Name())
		if h, ok := m.Handler(handlerName); ok {
			result := action.Call([]reflect.Value{h, reflect.ValueOf(ctx)})
			if len(result) > 0 {
				return result[0].Interface()
			}
		}
	}
	return nil
}

func (m Muxer) actionName(action reflect.Value) string {
	rawName := runtime.FuncForPC(action.Pointer()).Name()
	lastIndex := LastIndex(rawName, ".")
	name := ToLower(rawName[lastIndex+1:])
	return name
}

func (m Muxer) handlerName(action reflect.Value) string {
	actionType := action.Type()
	handlerType := actionType.In(0)
	handlerName := ToLower(handlerType.Name())
	lastIndex := LastIndex(handlerName, "handler")
	return handlerName[:lastIndex]
}

func (m Muxer) getViewForAction(action reflect.Value) string {
	actionName := m.actionName(action)
	handlerName := m.handlerName(action)
	return filepath.Join(handlerName, actionName)
}

func (m *Muxer) handleAppError(err *AppError, w http.ResponseWriter, bw *bufio.Writer) {
	errView := fmt.Sprintf("error/%v", err.Status())
	log.Printf("Err: %s\n%s", err, err.Stack())
	w.WriteHeader(err.Status())
	if ViewExists(errView) {
		bw.WriteString(fmt.Sprintf("Err: %s\n%s", err, err.Stack()))
		//RenderView(errView, bw, err.Error())
	} else {
		bw.WriteString(err.Error())
	}
}

func (m *Muxer) handleRedirect(ctx *Context) {
	var path string
	switch t := ctx.redirect.(type) {
	case string:
		path = t
	default:
		var ok bool
		path, ok = m.ActionPath(ctx.redirect)
		if !ok {
			panic("Did not find path for redirection, action not registered")
		}
	}

	http.Redirect(
		ctx.Writer,
		ctx.Request,
		path,
		http.StatusTemporaryRedirect,
	)
}
