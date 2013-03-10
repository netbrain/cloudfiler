package handler

import (
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/web"
	"net/http"
)

func NewInitHandler(userController UserController, roleController RoleController) InitHandler {
	return InitHandler{
		userController: userController,
		roleController: roleController,
	}
}

type InitHandler struct {
	data           UserData
	userController UserController
	roleController RoleController
}

func (i InitHandler) Init(ctx *Context) interface{} {
	if count, _ := i.userController.Count(); count > 0 {
		return Error(http.StatusNotFound)
	}
	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		data := &i.data
		ctx.InjectData(data)
		if err := data.Validate(); err != nil {
			ctx.AddValidationError(err)
		} else {
			user, err := i.userController.Create(data.Email, data.Password)
			if err != nil {
				return Error(err)
			}
			role, err := i.roleController.Create("Admin")
			if err != nil {
				return Error(err)
			}

			err = i.roleController.AddUser(role, user)
			if err != nil {
				return Error(err)
			}

			//TODO add flash message
			ctx.Redirect(AuthHandler.Login)
		}

	}
	return nil
}
