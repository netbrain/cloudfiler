package handler

import (
	"fmt"
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/web"
	"net/http"
)

type RoleData struct {
	Id   int
	Name string
}

type RoleHandler struct {
	controller RoleController
	data       RoleData
}

func NewRoleHandler(c RoleController) RoleHandler {
	return RoleHandler{
		controller: c,
	}
}

func (h RoleHandler) List(ctx *Context) interface{} {
	ctx.SetHeader(
		"Cache-Control",
		"no-cache, max-age=0, must-revalidate, no-store",
	)
	roles, err := h.controller.Roles()
	if err != nil {
		Error(fmt.Errorf("Could not fetch role list due to error: %v", err))
	}
	return roles
}

func (h RoleHandler) Create(ctx *Context) interface{} {
	ctrl := h.controller
	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		data := &h.data
		ctx.InjectData(data)

		if role, _ := ctrl.RoleByName(data.Name); role != nil {
			//TODO
			//ctx.AddFlashMessage("Role already registered")
		}
		if err := ctrl.Create(data.Name); err != nil {
			return Error(err)
		}

		ctx.Redirect(RoleHandler.List)

	}
	return nil
}

func (h RoleHandler) Retrieve(ctx *Context) interface{} {
	data := &h.data
	ctx.InjectData(data)

	role, err := h.controller.Role(data.Id)
	if err != nil {
		return Error(err)
	}

	if role == nil {
		return Error(http.StatusNotFound)
	}

	return role
}

func (h RoleHandler) Update(ctx *Context) interface{} {
	data := &h.data
	ctx.InjectData(data)

	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {

		if err := h.controller.Update(data.Id, data.Name); err != nil {
			return Error(err)
		}
		ctx.Redirect(RoleHandler.List)
		return nil

	}
	return h.Retrieve(ctx)

}

func (h RoleHandler) Delete(ctx *Context) interface{} {
	data := &h.data
	ctx.InjectData(data)

	if err := h.controller.Delete(data.Id); err != nil {
		return Error(err)
	}

	ctx.Redirect(RoleHandler.List)
	return nil
}
