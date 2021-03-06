package auth

import (
	"bufio"
	"errors"
	"fmt"
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	. "github.com/netbrain/cloudfiler/app/web/testing"
	"net/http"
	"strings"
	"testing"
)

var authenticator Authenticator
var userRepo UserRepositoryMem
var roleRepo RoleRepositoryMem

func InitAuthenticatorTest() {
	userRepo = NewUserRepository()
	roleRepo = NewRoleRepository()
	roleRepo.Store(&Role{
		Name: "Admin",
	})
	roleRepo.Store(&Role{
		Name: "User",
	})
	authenticator = NewAuthenticator(userRepo, roleRepo, "", "")
}

func user() (*User, string) {
	user := &User{
		Email: "test@test.test",
	}
	password := "testpass"
	user.SetPassword("testpass")
	userRepo.Store(user)
	return user, password
}

func authorize(email, password string) (*http.Cookie, error) {
	ctx, w := CreateReqContext("GET", "/protected", nil)
	ok := authenticator.Authorize(email, password, ctx.Writer, ctx.Request)
	if !ok {
		return nil, errors.New("Failed to login")
	}
	if _, ok := w.HeaderMap["Set-Cookie"]; !ok {
		return nil, errors.New("Set-Cookie doesn't exist")
	}

	if len(w.HeaderMap["Set-Cookie"]) == 0 {
		return nil, errors.New("Expected cookie val")
	}
	cookie, err := parseCookie(w.HeaderMap["Set-Cookie"][0])
	return cookie, err
}

func parseCookie(cookie string) (*http.Cookie, error) {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", cookie))))
	if err != nil {
		return nil, err
	}
	cookies := req.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("No cookies!")
	}
	return cookies[0], nil
}

func TestAuthorize(t *testing.T) {
	InitAuthenticatorTest()
	user, password := user()
	if _, err := authorize(user.Email, password); err != nil {
		t.Fatal(err)
	}
}

func TestAuthorizeWithInvalidCredentials(t *testing.T) {
	InitAuthenticatorTest()
	if _, err := authorize("", ""); err == nil {
		t.Fatal("Should have failed")
	}
}

func TestIsAuthorizedWhenZeroUsers(t *testing.T) {
	InitAuthenticatorTest()
	ctx, _ := CreateReqContext("GET", "/", nil)
	if !authenticator.IsAuthorized(ctx.Request) {
		t.Fatal("Should not have failed here, Zero users in system should yield unlimited access")
	}
}

func TestIsAuthorizedWhenNoAuthProcessIsCompleted(t *testing.T) {
	InitAuthenticatorTest()
	user()
	ctx, _ := CreateReqContext("GET", "/", nil)
	if authenticator.IsAuthorized(ctx.Request) {
		t.Fatal("Should have failed here")
	}
}

func TestIsAuthorizedWhenValidAuthIsInPlace(t *testing.T) {
	InitAuthenticatorTest()
	user, password := user()
	cookie, err := authorize(user.Email, password)
	if err != nil {
		t.Fatal(err)
	}
	ctx, _ := CreateReqContext("GET", "/protected", nil)
	ctx.Request.AddCookie(cookie)

	if !authenticator.IsAuthorized(ctx.Request) {
		t.Fatal("Should be authenticated")
	}

	if authenticator.AuthorizedUserID(ctx.Request) != user.ID {
		t.Fatal("ID doesn't match")
	}
}

func TestUnauthorize(t *testing.T) {
	InitAuthenticatorTest()
	user, password := user()
	cookie, err := authorize(user.Email, password)
	if err != nil {
		t.Fatal(err)
	}

	ctx, _ := CreateReqContext("GET", "/protected", nil)
	ctx.Request.AddCookie(cookie)

	if !authenticator.IsAuthorized(ctx.Request) {
		t.Fatal("Should be authenticated")
	}

	if authenticator.AuthorizedUserID(ctx.Request) != user.ID {
		t.Fatal("ID doesn't match")
	}

	authenticator.Unauthorize(ctx.Writer, ctx.Request)

	if authenticator.IsAuthorized(ctx.Request) {
		t.Fatal("Should not be authenticated")
	}

	if authenticator.AuthorizedUserID(ctx.Request) == user.ID {
		t.Fatal("ID matches!")
	}
}

func TestAdminProtectedResourceWhereUserDontHaveAdminRights(t *testing.T) {
	InitAuthenticatorTest()
	authenticator.SetRequiredPrivileges("/admin", "Admin")

	user, password := user()
	cookie, err := authorize(user.Email, password)
	if err != nil {
		t.Fatal(err)
	}

	ctx, _ := CreateReqContext("GET", "/admin", nil)
	ctx.Request.AddCookie(cookie)

	if authenticator.Handle(ctx.Writer, ctx.Request) {
		t.Fatal("User doesnt have admin rights!")
	}

}

func TestAdminProtectedResourceWhereUserHasAdminRights(t *testing.T) {
	InitAuthenticatorTest()
	authenticator.SetRequiredPrivileges("/admin", "Admin")

	user, password := user()

	role, _ := roleRepo.FindByName("Admin")
	role.Users = append(role.Users, *user)
	roleRepo.Store(role)

	cookie, err := authorize(user.Email, password)
	if err != nil {
		t.Fatal(err)
	}

	ctx, _ := CreateReqContext("GET", "/admin", nil)
	ctx.Request.AddCookie(cookie)

	if !authenticator.Handle(ctx.Writer, ctx.Request) {
		t.Fatal("User has admin rights!")
	}

}
