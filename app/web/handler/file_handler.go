package handler

import (
	. "github.com/netbrain/cloudfiler/app/web"
)

type FileHandler struct{}

func NewFileHandler() FileHandler {
	return FileHandler{}
}

func (h FileHandler) List(ctx *Context) interface{} {
	ctx.SetHeader(
		"Cache-Control",
		"no-cache, max-age=0, must-revalidate, no-store",
	)
	return nil
}
