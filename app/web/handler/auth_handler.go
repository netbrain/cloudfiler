package handler

import (
	. "github.com/netbrain/cloudfiler/app/web"
	. "github.com/netbrain/cloudfiler/app/web/auth"
)

type AuthData struct {
	Username string
	Password string
}

type AuthHandler struct {
	authenticator *Authenticator
	data          AuthData
}

func NewAuthHandler(authenticator *Authenticator) AuthHandler {
	return AuthHandler{
		authenticator: authenticator,
	}
}

func redirectToLandingPage(ctx *Context) {
	ctx.Redirect("/")
}

func (h AuthHandler) Login(ctx *Context) interface{} {
	if h.authenticator.IsAuthorized(ctx.Request) {
		redirectToLandingPage(ctx)
		return nil
	}

	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		data := &h.data
		ctx.InjectData(data)
		if ok := h.authenticator.Authorize(data.Username, data.Password, ctx.Writer, ctx.Request); ok {
			redirectToLandingPage(ctx)
		} else {
			//TODO add flash message
		}
	}
	return nil
}

func (h AuthHandler) Logout(ctx *Context) interface{} {
	h.authenticator.Unauthorize(ctx.Writer, ctx.Request)
	redirectToLandingPage(ctx)
	return nil
}
