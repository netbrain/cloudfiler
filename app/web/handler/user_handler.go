package handler

import (
	"fmt"
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/web"
	"net/http"
)

type UserData struct {
	Id            int
	Email         string
	Password      string
	PasswordAgain string `name:"password-again"`
}

func (data *UserData) Validate() error {
	if data.Password != data.PasswordAgain {
		return Error("Password did not match")
	}
	return nil
}

type UserHandler struct {
	controller UserController
	data       UserData
}

func NewUserHandler(controller UserController) UserHandler {
	handler := UserHandler{
		controller: controller,
		data:       UserData{},
	}
	return handler
}

func (handler UserHandler) List(ctx *Context) interface{} {
	ctx.SetHeader(
		"Cache-Control",
		"no-cache, max-age=0, must-revalidate, no-store",
	)
	users, err := handler.controller.Users()
	if err != nil {
		Error(fmt.Errorf("Could not fetch user list due to error: %v", err))
	}
	return users
}

func (handler UserHandler) Create(ctx *Context) interface{} {
	ctrl := handler.controller
	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		data := &handler.data
		ctx.InjectData(data)

		//Do custom validation
		if err := data.Validate(); err != nil {
			ctx.AddValidationError(err)
		} else {
			if user, _ := ctrl.UserByEmail(data.Email); user != nil {
				ctx.AddFlash("Email already registered")
			}
			if _, err := ctrl.Create(data.Email, data.Password); err != nil {
				return Error(err)
			}

			ctx.Redirect(UserHandler.List)
		}
	}
	return nil
}

func (handler UserHandler) Retrieve(ctx *Context) interface{} {
	data := &handler.data
	ctx.InjectData(data)

	user, err := handler.controller.User(data.Id)
	if err != nil {
		return Error(err)
	}

	if user == nil {
		return Error(http.StatusNotFound)
	}

	return user
}

func (handler UserHandler) Update(ctx *Context) interface{} {
	data := &handler.data
	ctx.InjectData(data)

	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		if err := data.Validate(); err != nil {
			ctx.AddValidationError(err)
		} else {
			if err := handler.controller.Update(data.Id, data.Email, data.Password); err != nil {
				return Error(err)
			}
			ctx.Redirect(UserHandler.List)
			return nil
		}
	}
	return handler.Retrieve(ctx)

}
