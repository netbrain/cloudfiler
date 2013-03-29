package handler

import (
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	. "github.com/netbrain/cloudfiler/app/web"
	. "github.com/netbrain/cloudfiler/app/web/auth"
	"io"
	"net/http"
	"time"
)

type FileHandler struct {
	fileController FileController
	authenticator  Authenticator
	data           struct {
		//File *multipart.FileHeader
		Id int
	}
}

func NewFileHandler(authenticator Authenticator, fileController FileController) FileHandler {
	return FileHandler{
		fileController: fileController,
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

	return data
}

func (h FileHandler) Upload(ctx *Context) interface{} {
	if ctx.Method() == "POST" && !ctx.HasValidationErrors() {
		//TODO allow file to be injected to data obj
		user, _ := h.authenticator.AuthorizedUser(ctx.Request)
		for _, fh := range ctx.Request.MultipartForm.File["file"] {
			fdata := new(FileDataMem) //TODO get FileData impl from configuration
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
	data := &h.data
	ctx.InjectData(data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	}

	if file == nil {
		return Error(http.StatusNotFound)
	}

	ctx.SetRawResponse(true)
	ctx.SetHeader("Content-Disposition", "attachment; filename="+file.Name)
	http.ServeContent(ctx.Writer, ctx.Request, file.Name, time.Time{}, file.Data)

	return nil
}

func (h FileHandler) Retrieve(ctx *Context) interface{} {
	data := &h.data
	ctx.InjectData(data)

	file, err := h.fileController.File(data.Id)
	if err != nil {
		return Error(err)
	}

	if file == nil {
		return Error(http.StatusNotFound)
	}

	return file
}
