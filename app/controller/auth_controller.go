package controller

import (
	"github.com/gorilla/sessions"
	. "github.com/netbrain/cloudfiler/app/conf"
	. "github.com/netbrain/cloudfiler/app/repository"
	"net/http"
)

const (
	AUTH_COOKIE_NAME = "auth"
	AUTH_ID          = "authId"
)

type AuthController struct {
	sessionStore   *sessions.CookieStore
	userRepository UserRepository
}

func NewAuthController(userRepository UserRepository) AuthController {
	c := AuthController{
		userRepository: userRepository,
		sessionStore:   sessions.NewCookieStore([]byte(Config.CookieStoreSecret)),
	}
	return c
}

func (c *AuthController) IsAuthenticated(r *http.Request) bool {
	session, _ := c.sessionStore.Get(r, AUTH_COOKIE_NAME)
	if session.IsNew {
		return false
	}

	if _, ok := session.Values[AUTH_ID]; ok {
		return true
	}

	return false
}

func (c *AuthController) Authenticate(email, password string, r *http.Request, w http.ResponseWriter) bool {
	user, _ := c.userRepository.FindByEmail(email)
	if user != nil && user.PasswordEquals(password) {
		session, _ := c.sessionStore.Get(r, AUTH_COOKIE_NAME)
		session.Values[AUTH_ID] = user.ID
		session.Save(r, w)
		return true
	}
	return false
}
