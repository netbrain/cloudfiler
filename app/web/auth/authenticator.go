package auth

import (
	"github.com/gorilla/sessions"
	. "github.com/netbrain/cloudfiler/app/conf"
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository"
	"log"
	"net/http"
)

const (
	AUTH_COOKIE_NAME = "auth"
	AUTH_ID          = "authId"
)

type Authenticator struct {
	sessionStore       *sessions.CookieStore
	userRepository     UserRepository
	roleRepository     RoleRepository
	loginUrl           string
	entryUrl           string
	requiredPrivileges map[string][]string
}

func NewAuthenticator(userRepository UserRepository, roleRepository RoleRepository, loginUrl, entryUrl string) Authenticator {
	a := Authenticator{
		userRepository:     userRepository,
		roleRepository:     roleRepository,
		sessionStore:       sessions.NewCookieStore(Config.CookieStoreAuthenticationKey, Config.CookieStoreEncryptionKey),
		loginUrl:           loginUrl,
		entryUrl:           entryUrl,
		requiredPrivileges: make(map[string][]string),
	}
	a.sessionStore.Options = &sessions.Options{
		Path: "/",
	}
	return a
}

func (a *Authenticator) getSession(r *http.Request) (*sessions.Session, error) {
	return a.sessionStore.Get(r, AUTH_COOKIE_NAME)
}

func (a *Authenticator) userHasRole(user *User, roleName string) bool {
	hasRole := false
	if user != nil {
		role, _ := a.roleRepository.FindByName(roleName)
		if len(role.Users) > 0 {
			for _, otherUser := range role.Users {
				if user.Equals(otherUser) {
					hasRole = true
					break
				}
			}
		}
	}
	return hasRole
}

func (a *Authenticator) userHasRoles(user *User, roleNames ...string) bool {
	hasRoles := true
	for _, roleName := range roleNames {
		if !a.userHasRole(user, roleName) {
			hasRoles = false
			break
		}
	}
	return hasRoles
}

func (a *Authenticator) SetRequiredPrivileges(path string, roles ...string) {
	a.requiredPrivileges[path] = roles
}

func (a *Authenticator) IsAuthorized(r *http.Request) bool {
	userCount, _ := a.userRepository.Count()
	if userCount == 0 {
		return true
	}

	session, _ := a.getSession(r)
	sessVal, authorized := session.Values[AUTH_ID]

	if authorized && a.requiredPrivileges != nil {
		authId := sessVal.(int)
		requiredRoles, exists := a.requiredPrivileges[r.URL.Path]
		if exists {
			user, _ := a.userRepository.FindById(authId)
			authorized = a.userHasRoles(user, requiredRoles...)
		}
	}

	return authorized
}

func (a *Authenticator) Authorize(email, password string, w http.ResponseWriter, r *http.Request) bool {
	user, _ := a.userRepository.FindByEmail(email)
	if user != nil && user.PasswordEquals(password) {
		session, _ := a.getSession(r)
		session.Options.MaxAge = 30 * 60
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

func (a Authenticator) Handle(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == a.loginUrl {
		return true
	}

	if a.IsAuthorized(r) {
		return true
	}

	//TODO add flash message
	http.Redirect(w, r, a.loginUrl, http.StatusTemporaryRedirect)
	return false
}
