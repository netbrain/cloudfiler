package handler

import (
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	"github.com/netbrain/cloudfiler/app/web"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var repo UserRepositoryMem
var controller UserController
var handler UserHandler

func initUserHandlerTest() {
	repo = NewUserRepository()
	controller = NewUserController(repo)
	handler = NewUserHandler(controller)
}

func CreateReqContext(method, path string, parameters map[string][]string) (*web.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, path, nil)
	switch method {
	case "GET":
		for key, slice := range parameters {
			for _, val := range slice {
				r.URL.Query().Add(key, val)
			}
		}
	case "POST":
		r.Form = parameters
	default:
		panic("unknown method")
	}

	ctx := web.NewContext(w, r)
	return ctx, w
}

func TestCreateReqContext(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("POST", "/test/path", nil)
	if ctx.Method() != "POST" {
		t.Fatal("Expected POST as method")
	}

	p := ctx.Request.URL.Path
	if p != "/test/path" {
		t.Fatal("Expected /test/path")
	}
}

func TestGetList(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("GET", "/user/list", nil)
	result := handler.List(ctx)
	if result == nil {
		t.Fatalf("response returned %v", result)
	}

}
func TestGetCreatePage(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("GET", "/user/create", nil)
	result := handler.Create(ctx)
	if result != nil {
		t.Fatal("Expected nil")
	}
}

func TestCreateUser(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("POST", "/user/create", map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"testpassword"},
		"password-again": []string{"testpassword"},
	})
	handler.Create(ctx)

	if ctx.HasValidationErrors() {
		t.Log(ctx.ValidationErrors)
		t.Fatal("Did not expect any errors")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

	result, _ := repo.All()
	if len(result) != 1 {
		t.Fatal("Expected user created")
	}
}

func TestCreateUserWithNoParameters(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("POST", "/user/create", nil)

	handler.Create(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}
}

func TestCreateUserWithInvalidEmail(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("POST", "/user/create", map[string][]string{
		"email":          []string{"test[at]test.test"},
		"password":       []string{"testpassword"},
		"password-again": []string{"testpassword"},
	})

	handler.Create(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}

	if ctx.ValidationErrors["email"] == nil {
		t.Fatal("Expected failure on email field")
	}
}

func TestCreateUserWithSmallPasswordLength(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("POST", "/user/create", map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"pass"},
		"password-again": []string{"pass"},
	})

	handler.Create(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}

	if ctx.ValidationErrors["password"] == nil {
		t.Fatal("Expected failure on password field")
	}
}

func TestCreateUserWithLargePasswordLength(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("POST", "/user/create", map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"my-very-very-very-long-password"},
		"password-again": []string{"my-very-very-very-long-password"},
	})

	handler.Create(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}

	if ctx.ValidationErrors["password"] == nil {
		t.Fatal("Expected failure on password field")
	}
}

func TestCreateUserWithMismatchingPasswordEntries(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("POST", "/user/create", map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"myvalidpass"},
		"password-again": []string{"mymismatchingpass"},
	})

	handler.Create(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}

	if ctx.ValidationErrors[""] == nil {
		t.Fatal("Expected a general error")
	}
}

func TestCreateUserWhereEmailAlreadyExist(t *testing.T) {
	initUserHandlerTest()
	controller.Create("test@test.test", "mypassword")

	ctx, _ := CreateReqContext("POST", "/user/create", map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"myvalidpass"},
		"password-again": []string{"myvalidpass"},
	})

	result := handler.Create(ctx)

	switch result.(type) {
	case error:
		//Expects a application error
	default:
		t.Fatal("Unexpected return type")
	}
}

func TestRetrieveNonExistingUser(t *testing.T) {
	initUserHandlerTest()

	ctx, _ := CreateReqContext("GET", "/user/retrieve", map[string][]string{
		"id": []string{"123"},
	})

	result := handler.Retrieve(ctx)

	if err, ok := result.(*web.AppError); ok {
		if err.Status() != http.StatusNotFound {
			t.Fatal("Expected 404")
		}
	} else {
		t.Fatal("Expected app error")
	}
}

func TestRetrieveExistingUser(t *testing.T) {
	initUserHandlerTest()

	controller.Create("test@test.test", "password")
	user, _ := repo.FindByEmail("test@test.test")
	id := strconv.Itoa(user.ID)

	ctx, _ := CreateReqContext("GET", "/user/retrieve?id="+id, nil)

	result := handler.Retrieve(ctx)

	if !user.Equals(result) {
		t.Fatal("Expected equality")
	}
}

func TestUpdateUser(t *testing.T) {
	initUserHandlerTest()
	controller.Create("test@test.test", "password")
	user, _ := repo.FindByEmail("test@test.test")
	id := strconv.Itoa(user.ID)

	ctx, _ := CreateReqContext("POST", "/user/update", map[string][]string{
		"id":             []string{id},
		"email":          []string{"test@test.test"},
		"password":       []string{"testpassword"},
		"password-again": []string{"testpassword"},
	})
	handler.Update(ctx)

	if ctx.HasValidationErrors() {
		t.Log(ctx.ValidationErrors)
		t.Fatal("Did not expect any errors")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

	user, _ = repo.FindByEmail("test@test.test")
	if !user.PasswordEquals("testpassword") {
		t.Fatal("Expected a updated password")
	}
}

func TestUpdateUserWithNoParameters(t *testing.T) {
	initUserHandlerTest()
	controller.Create("test@test.test", "password")
	user, _ := repo.FindByEmail("test@test.test")
	id := strconv.Itoa(user.ID)
	ctx, _ := CreateReqContext("POST", "/user/update?id="+id, nil)

	handler.Update(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}
}

func TestUpdateUserWithInvalidEmail(t *testing.T) {
	initUserHandlerTest()
	controller.Create("test@test.test", "password")
	user, _ := repo.FindByEmail("test@test.test")
	id := strconv.Itoa(user.ID)
	ctx, _ := CreateReqContext("POST", "/user/update", map[string][]string{
		"id":             []string{id},
		"email":          []string{"test[at]test.test"},
		"password":       []string{"testpassword"},
		"password-again": []string{"testpassword"},
	})

	handler.Update(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}

	if ctx.ValidationErrors["email"] == nil {
		t.Fatal("Expected failure on email field")
	}
}

func TestUpdateUserWithSmallPasswordLength(t *testing.T) {
	initUserHandlerTest()
	controller.Create("test@test.test", "password")
	user, _ := repo.FindByEmail("test@test.test")
	id := strconv.Itoa(user.ID)
	ctx, _ := CreateReqContext("POST", "/user/update?id="+id, map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"pass"},
		"password-again": []string{"pass"},
	})

	handler.Update(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}

	if ctx.ValidationErrors["password"] == nil {
		t.Fatal("Expected failure on password field")
	}
}

func TestUpdateUserWithLargePasswordLength(t *testing.T) {
	initUserHandlerTest()
	controller.Create("test@test.test", "password")
	user, _ := repo.FindByEmail("test@test.test")
	id := strconv.Itoa(user.ID)
	ctx, _ := CreateReqContext("POST", "/user/update?id="+id, map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"my-very-very-very-long-password"},
		"password-again": []string{"my-very-very-very-long-password"},
	})

	handler.Update(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}

	if ctx.ValidationErrors["password"] == nil {
		t.Fatal("Expected failure on password field")
	}
}

func TestUpdateUserWithMismatchingPasswordEntries(t *testing.T) {
	initUserHandlerTest()
	controller.Create("test@test.test", "password")
	user, _ := repo.FindByEmail("test@test.test")
	id := strconv.Itoa(user.ID)
	ctx, _ := CreateReqContext("POST", "/user/update?id="+id, map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"myvalidpass"},
		"password-again": []string{"mymismatchingpass"},
	})

	handler.Update(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}

	if ctx.ValidationErrors[""] == nil {
		t.Fatal("Expected a general error")
	}
}

func TestDeleteNonExistingUser(t *testing.T) {
	initUserHandlerTest()
	ctx, _ := CreateReqContext("GET", "/user/delete?id=123", nil)

	result := handler.Delete(ctx)

	if result != nil {
		t.Fatal("expected nil")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

}

func TestDeleteUser(t *testing.T) {
	initUserHandlerTest()
	controller.Create("test@test.test", "password")
	user, _ := repo.FindByEmail("test@test.test")
	id := strconv.Itoa(user.ID)
	ctx, _ := CreateReqContext("GET", "/user/delete?id="+id, nil)

	result := handler.Delete(ctx)

	if result != nil {
		t.Fatal("expected nil")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}
}
