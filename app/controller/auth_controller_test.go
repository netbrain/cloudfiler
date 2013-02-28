package controller

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	"net/http"
	"net/http/httptest"
	"testing"
)

var authController AuthController

func InitAuthControllerTest() {
	userRepo = NewUserRepository()
	authController = NewAuthController(userRepo)
}

func TestAuthenticate(t *testing.T) {
	InitAuthControllerTest()
	r, _ := http.NewRequest("GET", "/auth", nil)
	w := httptest.NewRecorder()
	user := &User{
		Email: "test@test.test",
	}
	user.SetPassword("testpass")
	userRepo.Store(user)
	ok := authController.Authenticate("test@test.test", "testpass", r, w)
	if !ok {
		t.Fatal("Failed to login")
	}

	if len(w.HeaderMap["Set-Cookie"]) == 0 {
		t.Fatal("Expected cookie")
	}
}

func TestAuthenticateWithInvalidCredentials(t *testing.T) {
	InitAuthControllerTest()
	r, _ := http.NewRequest("GET", "/auth", nil)
	w := httptest.NewRecorder()

	ok := authController.Authenticate("", "", r, w)
	if ok {
		t.Fatal("Should have failed to login")
	}
}

func TestIsAuthenticatedWhenNoAuthProcessIsCompleted(t *testing.T) {
	InitAuthControllerTest()
	r, _ := http.NewRequest("GET", "/", nil)
	if authController.IsAuthenticated(r) {
		t.Fatal("Should have failed here")
	}
}

func TestIsAuthenticatedWhenValidAuthIsInPlace(t *testing.T) {
	InitAuthControllerTest()
	r, _ := http.NewRequest("GET", "/auth", nil)
	w := httptest.NewRecorder()
	user := &User{
		Email: "test@test.test",
	}
	user.SetPassword("testpass")
	userRepo.Store(user)
	ok := authController.Authenticate("test@test.test", "testpass", r, w)
	if !ok {
		t.Fatal("Failed to login")
	}

	r, _ = http.NewRequest("GET", "/protected", nil)
	r.Header.Set("Cookie", w.HeaderMap["Set-Cookie"][0])

	if !authController.IsAuthenticated(r) {
		t.Fatal("Should be authenticated")
	}
}
