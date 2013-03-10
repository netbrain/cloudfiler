package handler

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	. "github.com/netbrain/cloudfiler/app/web/auth"
	. "github.com/netbrain/cloudfiler/app/web/testing"
	"testing"
)

var authenticator Authenticator
var authHandler AuthHandler

func initAuthHandlerTest() {
	userRepo = NewUserRepository()
	roleRepo = NewRoleRepository()
	authenticator = NewAuthenticator(userRepo, roleRepo, "", "")
	authHandler = NewAuthHandler(authenticator)
}

func TestGetLoginPage(t *testing.T) {
	initAuthHandlerTest()
	ctx, _ := CreateReqContext("GET", "/auth/login", nil)
	result := authHandler.Login(ctx)

	if result != nil {
		t.Fatal("Expected nil")
	}

	if ctx.IsRedirected() {
		t.Fatal("Should not redirect")
	}
}

func TestFailedLogin(t *testing.T) {
	initAuthHandlerTest()
	ctx, _ := CreateReqContext("POST", "/auth/login", map[string][]string{
		"username": []string{"nonexistant@email.test"},
		"password": []string{"testpass"},
	})
	result := authHandler.Login(ctx)

	if result != nil {
		t.Fatal("Expected nil")
	}

	if ctx.IsRedirected() {
		t.Fatal("Should not redirect")
	}
}

func TestSuccessfullLogin(t *testing.T) {
	initAuthHandlerTest()
	user := &User{
		Email: "testuser@test.test",
	}
	password := "testpass"
	user.SetPassword(password)
	userRepo.Store(user)

	ctx, _ := CreateReqContext("POST", "/auth/login", map[string][]string{
		"username": []string{user.Email},
		"password": []string{password},
	})
	result := authHandler.Login(ctx)

	if result != nil {
		t.Fatal("Expected nil")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Should redirect to landing page")
	}
}

func TestLogout(t *testing.T) {
	initAuthHandlerTest()
	user := &User{
		Email: "testuser@test.test",
	}
	password := "testpass"
	user.SetPassword(password)
	userRepo.Store(user)

	ctx, _ := CreateReqContext("POST", "/auth/login", map[string][]string{
		"username": []string{user.Email},
		"password": []string{password},
	})
	result := authHandler.Login(ctx)

	if result != nil {
		t.Fatal("Expected nil")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Should redirect to landing page")
	}

	ctx, _ = CreateReqContext("GET", "/auth/logout", nil)
	result = authHandler.Logout(ctx)

	if !ctx.IsRedirected() {
		t.Fatal("Should redirect to landing page")
	}
}
