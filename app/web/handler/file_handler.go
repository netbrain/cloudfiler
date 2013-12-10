package handler

import (
	"encoding/json"
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/entity"
	"github.com/netbrain/cloudfiler/app/repository/fs"
	. "github.com/netbrain/cloudfiler/app/web"
	. "github.com/netbrain/cloudfiler/app/web/auth"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type FileHandler struct {
	fileController FileController
	userController UserController
	roleController RoleController
	authenticator  Authenticator
}

func NewFileHandler(
	authenticator Authenticator,
	fileController FileController,
	userController UserController,
	roleController RoleController) FileHandler {

	return FileHandler{
		fileController: fileController,
		userController: userController,
		roleController: roleController,
		authenticator:  authenticator,
	}
}

func (h FileHandler) List(ctx *Context) interface{} {
	ctx.SetHeader(
		"Cache-Control",
		"no-cache, max-age=0, must-revalidate, no-store",
	)
	user, err := h.authenticator.AuthorizedUser(ctx.Request)

	if err != nil {
		return Error(err)
	}

	data, err := h.fileController.FilesWhereUserHasAccess(*user)
	if err != nil {
		return Error(err)
	}

	sort.Sort(ByUploaded(data))

	return data
}

func (h FileHandler) Upload(ctx *Context) interface{} {
	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		//TODO allow file to be injected to data obj
		user, _ := h.authenticator.AuthorizedUser(ctx.Request)
		for _, fh := range ctx.Request.MultipartForm.File["file"] {
			fdata := fs.NewFileData() //TODO get FileData impl from configuration
			file, err := fh.Open()
			if err != nil {
				return Error(err)
			}
			buffer := make([]byte, 1<<20) //1mb
			for {
				read, err := file.Read(buffer)

				if err != nil && err != io.EOF {
					return Error(err)
				}

				fdata.Write(buffer[:read])

				if read < len(buffer) {
					break
				}
			}

			file.Close()
			h.fileController.Create(fh.Filename, *user, fdata)
		}
		ctx.Redirect(FileHandler.List)
	}
	return nil
}

func (h FileHandler) Download(ctx *Context) interface{} {
	data := struct {
		Id int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	}

	if file == nil {
		return Error(http.StatusNotFound)
	}

	if h.hasAccess(file, ctx) {
		ctx.SetRawResponse(true)
		ctx.SetHeader("Content-Disposition", "attachment; filename="+file.Name)
		ctx.SetHeader("Content-Type", "application/octet-stream")
		ctx.SetHeader("Pragma", "no-cache")
		ctx.SetHeader("Expires", " 0")
		http.ServeContent(ctx.Writer, ctx.Request, file.Name, time.Time{}, file.Data)
		file.Data.Close()
	} else {
		//Used doesn't have access
		return Error(http.StatusForbidden)
	}
	return nil
}

func (h FileHandler) Retrieve(ctx *Context) interface{} {
	data := struct {
		Id int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	}

	if file == nil {
		return Error(http.StatusNotFound)
	}

	if h.hasAccess(file, ctx) {
		return file
	}
	return Error(http.StatusForbidden)
}

func (h FileHandler) AddUsers(ctx *Context) interface{} {
	data := struct {
		Uid []int
		Id  int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	} else if file == nil {
		return Error(http.StatusNotFound)
	}

	if h.hasAccess(file, ctx) {
		if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
			for _, uid := range data.Uid {
				user, err := h.userController.User(uid)

				if err != nil {
					return Error(err)
				} else if user == nil {
					continue
				}

				if err := h.fileController.GrantUserAccessToFile(*user, file); err != nil {
					return Error(err)
				}
			}
			ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
			return nil
		}

		users, _ := h.userController.Users()
		out := struct {
			File  *File
			Users []User
		}{
			File:  file,
			Users: users,
		}

		return out
	}
	ctx.AddFlash("Successfully added user.")
	ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
	return nil
}

func (h FileHandler) RemoveUsers(ctx *Context) interface{} {
	data := struct {
		Uid []int
		Id  int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	} else if file == nil {
		return Error(http.StatusNotFound)
	}

	if h.hasAccess(file, ctx) {
		for _, uid := range data.Uid {
			user, err := h.userController.User(uid)

			if err != nil {
				return Error(err)
			} else if user == nil {
				continue
			}

			if err := h.fileController.RevokeUserAccessToFile(*user, file); err != nil {
				return Error(err)
			}
		}
	}
	ctx.AddFlash("Successfully removed user.")
	ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
	return nil
}

func (h FileHandler) AddRoles(ctx *Context) interface{} {
	data := struct {
		Rid []int
		Id  int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	} else if file == nil {
		return Error(http.StatusNotFound)
	}

	if h.hasAccess(file, ctx) {
		if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
			for _, rid := range data.Rid {
				role, err := h.roleController.Role(rid)

				if err != nil {
					return Error(err)
				} else if role == nil {
					continue
				}

				if err := h.fileController.GrantRoleAccessToFile(*role, file); err != nil {
					return Error(err)
				}
			}
			ctx.AddFlash("Successfully added role.")
			ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
			return nil
		}

		roles, _ := h.roleController.Roles()
		out := struct {
			File  *File
			Roles []Role
		}{
			File:  file,
			Roles: roles,
		}

		return out
	}

	ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
	return nil
}

func (h FileHandler) RemoveRoles(ctx *Context) interface{} {
	data := struct {
		Rid []int
		Id  int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	} else if file == nil {
		return Error(http.StatusNotFound)
	}

	if h.hasAccess(file, ctx) {
		for _, rid := range data.Rid {
			role, err := h.roleController.Role(rid)

			if err != nil {
				return Error(err)
			} else if role == nil {
				continue
			}
			if err := h.fileController.RevokeRoleAccessToFile(*role, file); err != nil {
				return Error(err)
			}
		}
	}
	ctx.AddFlash("Successfully removed role.")
	ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
	return nil
}

func (h FileHandler) AddTags(ctx *Context) interface{} {
	data := struct {
		Tag []string
		Id  int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	} else if file == nil {
		return Error(http.StatusNotFound)
	}

	if !h.hasAccess(file, ctx) {
		return Error(http.StatusForbidden)
	}

	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		h.fileController.AddTags(file, data.Tag...)

		if !ctx.IsAjaxRequest() {
			ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
		}
	}

	ctx.SetRawResponse(true)
	return nil
}

func (h FileHandler) SetTags(ctx *Context) interface{} {
	data := struct {
		Tag []string
		Id  int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	} else if file == nil {
		return Error(http.StatusNotFound)
	}

	if !h.hasAccess(file, ctx) {
		return Error(http.StatusForbidden)
	}

	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		h.fileController.SetTags(file, data.Tag...)

		if !ctx.IsAjaxRequest() {
			ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
		}
	}

	ctx.SetRawResponse(true)
	return nil
}

func (h FileHandler) RemoveTags(ctx *Context) interface{} {
	data := struct {
		Tag []string
		Id  int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	} else if file == nil {
		return Error(http.StatusNotFound)
	}

	if !h.hasAccess(file, ctx) {
		return Error(http.StatusForbidden)
	}

	h.fileController.RemoveTags(file, data.Tag...)
	ctx.SetRawResponse(true)
	return nil
}

func (h FileHandler) Delete(ctx *Context) interface{} {
	data := struct {
		Id int
	}{}
	ctx.InjectData(&data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	} else if file == nil {
		return Error(http.StatusNotFound)
	}

	if h.hasAccess(file, ctx) {
		if err := h.fileController.Erase(file.ID); err != nil {
			return Error(err)
		}
	}
	ctx.AddFlash("Deleted file: " + file.Name)
	ctx.Redirect(FileHandler.List)
	return nil
}

func (h FileHandler) Search(ctx *Context) interface{} {
	data := struct {
		Query string
	}{}
	ctx.InjectData(&data)

	user, _ := h.authenticator.AuthorizedUser(ctx.Request)
	result, err := h.fileController.FileSearch(*user, data.Query)
	if err != nil {
		return Error(err)
	}

	return result
}

func (h FileHandler) Tags(ctx *Context) interface{} {
	tags := h.fileController.Tags()
	out, err := json.Marshal(tags)

	if err != nil {
		panic(err)
	}

	ctx.SetRawResponse(true)
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Writer.Write(out)
	return nil
}

func (h FileHandler) Update(ctx *Context) interface{} {
	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		data := struct {
			Id          int
			Description string
		}{}
		ctx.InjectData(&data)

		file, err := h.fileController.File(data.Id)
		if err != nil {
			return Error(err)
		} else if file == nil {
			return Error(http.StatusNotFound)
		}

		if h.hasAccess(file, ctx) {
			file.Description = data.Description
			if err := h.fileController.Update(file); err != nil {
				return Error(err)
			}
			ctx.AddFlash("File successfully updated.")
		}
		ctx.Redirect("/file/retrieve?id=" + strconv.Itoa(data.Id))
	}
	return nil
}

func (h FileHandler) hasAccess(file *File, ctx *Context) bool {
	user, err := h.authenticator.AuthorizedUser(ctx.Request)
	if err != nil {
		panic(err)
	} else if user == nil {
		panic("this should not happen, user should always be authenticated")
	}

	return h.fileController.UserHasAccess(*user, *file)
}
