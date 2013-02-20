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
	roleRepo := mem.NewRoleRepository()

	userController := controller.NewUserController(userRepo)
	roleController := controller.NewRoleController(roleRepo)

	userHandler := handler.NewUserHandler(userController)
	roleHandler := handler.NewRoleHandler(roleController)

	log.Println("Adding web handlers...")
	Muxer.AddHandler(userHandler)
	Muxer.AddHandler(roleHandler)
}

func initRoutes() {
	log.Println("Adding routing table...")

	//User
	Muxer.AddAction("/user/list", handler.UserHandler.List)
	Muxer.AddAction("/user/create", handler.UserHandler.Create)
	Muxer.AddAction("/user/retrieve", handler.UserHandler.Retrieve)
	Muxer.AddAction("/user/update", handler.UserHandler.Update)
	Muxer.AddAction("/user/delete", handler.UserHandler.Delete)

	//Role
	Muxer.AddAction("/role/list", handler.RoleHandler.List)
	Muxer.AddAction("/role/create", handler.RoleHandler.Create)
	Muxer.AddAction("/role/retrieve", handler.RoleHandler.Retrieve)
	Muxer.AddAction("/role/update", handler.RoleHandler.Update)
	Muxer.AddAction("/role/delete", handler.RoleHandler.Delete)

}
