package auth

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

type Authenticator struct {
	sessionStore   *sessions.CookieStore
	userRepository UserRepository
}

func NewAuthenticator(userRepository UserRepository) Authenticator {
	a := Authenticator{
		userRepository: userRepository,
		sessionStore:   sessions.NewCookieStore(Config.CookieStoreAuthenticationKey, Config.CookieStoreEncryptionKey),
	}
	a.sessionStore.Options = &sessions.Options{
		MaxAge: 30 * 60,
		Path:   "/",
	}
	return a
}

func (a *Authenticator) getSession(r *http.Request) (*sessions.Session, error) {
	return a.sessionStore.Get(r, AUTH_COOKIE_NAME)
}

func (a *Authenticator) IsAuthorized(r *http.Request) bool {
	allUsers, _ := a.userRepository.All()
	if len(allUsers) == 0 {
		return true
	}
	session, _ := a.getSession(r)
	_, authorized := session.Values[AUTH_ID]
	return authorized
}

func (a *Authenticator) Authorize(email, password string, w http.ResponseWriter, r *http.Request) bool {
	user, _ := a.userRepository.FindByEmail(email)
	if user != nil && user.PasswordEquals(password) {
		session, _ := a.getSession(r)

		session.Values[AUTH_ID] = user.ID
		if err := session.Save(r, w); err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func (a *Authenticator) Unauthorize(w http.ResponseWriter, r *http.Request) {
	session, _ := a.getSession(r)
	session.Options.MaxAge = -1
	delete(session.Values, AUTH_ID)
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
}

func (a *Authenticator) AuthorizedUserID(r *http.Request) int {
	session, _ := a.getSession(r)
	if val, ok := session.Values[AUTH_ID]; ok {
		if id, ok := val.(int); ok {
			return id
		}
	}
	return 0
}
