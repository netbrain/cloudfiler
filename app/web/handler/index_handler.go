package handler

import (
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/web"
	. "github.com/netbrain/cloudfiler/app/web/auth"
)

type IndexHandler struct {
	authenticator  Authenticator
	userController UserController
}

func NewIndexHandler(authenticator Authenticator, userController UserController) IndexHandler {
	return IndexHandler{
		authenticator:  authenticator,
		userController: userController,
	}
}

func (h IndexHandler) Index(ctx *Context) interface{} {

	if c, _ := h.userController.Count(); c > 0 {
		ctx.Redirect(FileHandler.List)
	} else {
		ctx.Redirect(InitHandler.Init)
	}

	return nil
}
