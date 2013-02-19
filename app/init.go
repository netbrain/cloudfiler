package app

import (
	"github.com/netbrain/cloudfiler/app/controller"
	"github.com/netbrain/cloudfiler/app/repository/mem"
	"github.com/netbrain/cloudfiler/app/web"
	"github.com/netbrain/cloudfiler/app/web/handlers"
	"log"
)

var Muxer = web.NewMuxer()

func init() {
	initApplication()
	initRoutes()
}

func initApplication() {
	log.Println("Initializing application dependencies...")
	userRepo := mem.NewUserRepository()
	userController := controller.NewUserController(userRepo)
	userHandler := handler.NewUserHandler(userController)

	log.Println("Adding web handlers...")
	Muxer.AddHandler(userHandler)
}

func initRoutes() {
	log.Println("Adding routing table...")
	Muxer.AddAction("/user/list", handler.UserHandler.List)
	Muxer.AddAction("/user/create", handler.UserHandler.Create)
	Muxer.AddAction("/user/retrieve", handler.UserHandler.Retrieve)
	Muxer.AddAction("/user/update", handler.UserHandler.Update)
	Muxer.AddAction("/user/delete", handler.UserHandler.Delete)
}
