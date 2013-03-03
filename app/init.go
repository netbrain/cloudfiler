package app

import (
	"github.com/netbrain/cloudfiler/app/controller"
	"github.com/netbrain/cloudfiler/app/repository/mem"
	"github.com/netbrain/cloudfiler/app/web"
	"github.com/netbrain/cloudfiler/app/web/auth"
	"github.com/netbrain/cloudfiler/app/web/handler"
	"github.com/netbrain/cloudfiler/app/web/interceptor"
	"log"
)

var muxer web.Muxer
var WebHandler = new(interceptor.InterceptorChain)

func init() {
	initApplication()
	initRoutes()
}

func initApplication() {
	log.Println("Initializing application dependencies...")
	userRepo := mem.NewUserRepository()
	roleRepo := mem.NewRoleRepository()
	authenticator := auth.NewAuthenticator(userRepo, "/auth/login", "/")
	WebHandler.AddInterceptor(authenticator)

	userController := controller.NewUserController(userRepo)
	roleController := controller.NewRoleController(roleRepo)

	userHandler := handler.NewUserHandler(userController)
	roleHandler := handler.NewRoleHandler(roleController)
	authHandler := handler.NewAuthHandler(authenticator)

	muxer = web.NewMuxer(authenticator)
	WebHandler.AddInterceptor(muxer)

	log.Println("Adding web handlers...")
	muxer.AddHandler(authHandler)
	muxer.AddHandler(userHandler)
	muxer.AddHandler(roleHandler)
}

func initRoutes() {
	log.Println("Adding routing table...")

	//Auth
	muxer.AddAction("/auth/login", handler.AuthHandler.Login)
	muxer.AddAction("/auth/logout", handler.AuthHandler.Logout)

	//User
	muxer.AddAction("/user/list", handler.UserHandler.List)
	muxer.AddAction("/user/create", handler.UserHandler.Create)
	muxer.AddAction("/user/retrieve", handler.UserHandler.Retrieve)
	muxer.AddAction("/user/update", handler.UserHandler.Update)
	muxer.AddAction("/user/delete", handler.UserHandler.Delete)

	//Role
	muxer.AddAction("/role/list", handler.RoleHandler.List)
	muxer.AddAction("/role/create", handler.RoleHandler.Create)
	muxer.AddAction("/role/retrieve", handler.RoleHandler.Retrieve)
	muxer.AddAction("/role/update", handler.RoleHandler.Update)
	muxer.AddAction("/role/delete", handler.RoleHandler.Delete)

}
