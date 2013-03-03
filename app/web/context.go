package web

import (
	"github.com/netbrain/cloudfiler/app/web/fvgo"
	"net/http"
)

type Context struct {
	ValidationErrors map[string][]error
	Data             interface{}
	Request          *http.Request
	Writer           http.ResponseWriter
	redirect         interface{}
	formValidator    *fvgo.FormValidator
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Request:       r,
		Writer:        w,
		formValidator: FormValidator,
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
	ctx.injectData(ctx.Request.URL.Query(), dataObject)
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
