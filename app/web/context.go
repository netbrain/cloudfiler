package web

import (
	"github.com/netbrain/cloudfiler/app/web/fvgo"
	"github.com/netbrain/cloudfiler/app/web/session"
	"net/http"
)

type Context struct {
	ValidationErrors map[string][]error
	Data             interface{}
	Request          *http.Request
	Writer           http.ResponseWriter
	redirect         interface{}
	rawResponse      bool
	formValidator    *fvgo.FormValidator
	session          *session.Session
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Request:       r,
		Writer:        w,
		formValidator: FormValidator,
		rawResponse:   false,
		session:       session.NewSession(w, r),
	}
	_, ctx.ValidationErrors = ctx.formValidator.ValidateRequestData(ctx.Request)
	return ctx
}

func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) AddHeader(key, value string) {
	ctx.Writer.Header().Add(key, value)
}

func (ctx *Context) DelHeader(key string) {
	ctx.Writer.Header().Del(key)
}

func (ctx *Context) GetHeader(key string) string {
	return ctx.Writer.Header().Get(key)
}

func (ctx *Context) GetRequestHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) Params(name string) string {
	return ctx.Request.FormValue(name)
}

func (ctx *Context) Redirect(handlerOrUrl interface{}) {
	ctx.redirect = handlerOrUrl
}

func (ctx *Context) IsRedirected() bool {
	return ctx.redirect != nil
}

func (ctx *Context) InjectData(dataObject interface{}) {
	ctx.Request.ParseForm()
	ctx.injectData(ctx.Request.Form, dataObject)
}

func (ctx *Context) injectData(inputData map[string][]string, dataObject interface{}) {
	fdi := NewDataInjector(inputData, dataObject)
	if fdi != nil {
		fdi.Inject()
	}
}

func (ctx *Context) AddValidationError(errs ...error) {
	ctx.AddFieldValidationError("", errs...)
}

func (ctx *Context) AddFieldValidationError(fieldName string, errs ...error) {
	if ctx.ValidationErrors == nil {
		ctx.ValidationErrors = make(map[string][]error)
	}
	if _, present := ctx.ValidationErrors[fieldName]; !present {
		ctx.ValidationErrors[fieldName] = make([]error, 0)
	}
	existing := ctx.ValidationErrors[fieldName]
	ctx.ValidationErrors[fieldName] = append(existing, errs...)
}

func (ctx *Context) HasValidationErrors() bool {
	return len(ctx.ValidationErrors) > 0
}

func (ctx *Context) SetRawResponse(b bool) {
	ctx.rawResponse = b
}

func (ctx *Context) IsAjaxRequest() bool {
	return ctx.GetRequestHeader("X-Requested-With") == "XMLHttpRequest"
}

func (ctx *Context) Session() *session.Session {
	return ctx.session
}

func (ctx *Context) Flash() interface{} {
	if flashes := ctx.session.Flash(); len(flashes) > 0 {
		return flashes
	}
	return nil
}

func (ctx *Context) AddFlash(v interface{}) {
	ctx.session.AddFlash(v)
}
