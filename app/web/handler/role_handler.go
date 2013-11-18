package handler

import (
	"fmt"
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/web"
	"net/http"
	"strconv"
)

type RoleData struct {
	Id   int
	Name string
}

type RoleHandler struct {
	roleController RoleController
	userController UserController
	data           RoleData
}

func NewRoleHandler(roleController RoleController, userController UserController) RoleHandler {
	return RoleHandler{
		roleController: roleController,
		userController: userController,
	}
}

func (h RoleHandler) List(ctx *Context) interface{} {
	ctx.SetHeader(
		"Cache-Control",
		"no-cache, max-age=0, must-revalidate, no-store",
	)
	roles, err := h.roleController.Roles()
	if err != nil {
		Error(fmt.Errorf("Could not fetch role list due to error: %v", err))
	}
	return roles
}

func (h RoleHandler) Create(ctx *Context) interface{} {
	ctrl := h.roleController
	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		data := &h.data
		ctx.InjectData(data)

		if role, _ := ctrl.RoleByName(data.Name); role != nil {
			ctx.AddFlash("Role already registered")
		}
		if _, err := ctrl.Create(data.Name); err != nil {
			return Error(err)
		}

		ctx.Redirect(RoleHandler.List)

	}
	return nil
}

func (h RoleHandler) Retrieve(ctx *Context) interface{} {
	data := &h.data
	ctx.InjectData(data)

	role, err := h.roleController.Role(data.Id)
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

		if err := h.roleController.Update(data.Id, data.Name); err != nil {
			return Error(err)
		}
		ctx.AddFlash("Role updated successfully.")
		//TODO make it possible to add parameters to redirect actions
		// e.g ctx.Redirect(RoleHandler.Retrieve,data.Id)
		ctx.Redirect("/role/retrieve?id=" + strconv.Itoa(data.Id))
		return nil

	}
	return h.Retrieve(ctx)
}

func (h RoleHandler) Delete(ctx *Context) interface{} {
	data := &h.data
	ctx.InjectData(data)

	if err := h.roleController.Delete(data.Id); err != nil {
		return Error(err)
	}

	ctx.AddFlash(fmt.Sprintf("Deleted role with id: %v", data.Id))
	ctx.Redirect(RoleHandler.List)
	return nil
}

func (h RoleHandler) AddUser(ctx *Context) interface{} {
	data := struct {
		Uid []int
		Id  int
	}{}
	ctx.InjectData(&data)

	role, err := h.roleController.Role(data.Id)
	if err != nil {
		return Error(err)
	} else if role == nil {
		return Error(http.StatusNotFound)
	}

	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		for _, uid := range data.Uid {
			user, err := h.userController.User(uid)

			if err != nil {
				return Error(err)
			} else if user == nil {
				continue
			}

			if err := h.roleController.AddUser(role, user); err != nil {
				return Error(err)
			}
		}
		//TODO make it possible to add parameters to redirect actions
		// e.g ctx.Redirect(RoleHandler.Retrieve,data.Id)
		ctx.AddFlash("User added.")
		ctx.Redirect("/role/retrieve?id=" + strconv.Itoa(data.Id))
		return nil
	}

	users, _ := h.userController.Users()
	out := struct {
		Role  *Role
		Users []User
	}{
		Role:  role,
		Users: users,
	}

	return out
}

func (h RoleHandler) RemoveUser(ctx *Context) interface{} {
	data := struct {
		Uid int
		Id  int
	}{}
	ctx.InjectData(&data)

	user, err := h.userController.User(data.Uid)
	if err != nil {
		return Error(err)
	}

	role, err := h.roleController.Role(data.Id)
	if err != nil {
		return Error(err)
	}

	if role == nil || user == nil {
		return Error(http.StatusNotFound)
	}
	if !(role.Name == "Admin" && len(role.Users) == 1) {
		if err := h.roleController.RemoveUser(role, user); err != nil {
			return Error(err)
		}
	} else {
		ctx.AddFlash("Cannot remove the only Admin user in the application")
		ctx.Redirect(ctx.GetRequestHeader("Referer"))
		return nil
	}
	ctx.AddFlash("User removed.")
	ctx.Redirect(ctx.GetRequestHeader("Referer"))
	return nil
}
