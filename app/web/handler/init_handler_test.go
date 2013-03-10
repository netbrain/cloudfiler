package handler

import (
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	. "github.com/netbrain/cloudfiler/app/web"
	. "github.com/netbrain/cloudfiler/app/web/testing"
	"net/http"
	//"strconv"
	"testing"
)

var initHandler InitHandler

func initInitHandlerTest() {
	userRepo = NewUserRepository()
	userController = NewUserController(userRepo)
	roleRepo = NewRoleRepository()
	roleController = NewRoleController(roleRepo, userRepo)
	initHandler = NewInitHandler(userController, roleController)
}

func TestGetInitWhenNoUsers(t *testing.T) {
	initInitHandlerTest()
	ctx, _ := CreateReqContext("GET", "/init", nil)
	result := initHandler.Init(ctx)
	if result != nil {
		t.Fatalf("response returned %v", result)
	}
}

func TestGetInitWhenOneOrMoreUsers(t *testing.T) {
	initInitHandlerTest()
	userController.Create("test@test.test", "password")
	ctx, _ := CreateReqContext("GET", "/init", nil)
	result := initHandler.Init(ctx)
	if err, ok := result.(*AppError); !ok {
		t.Fatalf("response returned %v", result)
	} else if err.Status() != http.StatusNotFound {
		t.Fatal("Expected 404")
	}
}

func TestInitCreateUser(t *testing.T) {
	initInitHandlerTest()
	ctx, _ := CreateReqContext("POST", "/init", map[string][]string{
		"email":          []string{"test@test.test"},
		"password":       []string{"testpasswd"},
		"password-again": []string{"testpasswd"},
	})
	result := initHandler.Init(ctx)
	if result != nil {
		t.Fatalf("response returned %v", result)
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

	if role, err := roleRepo.FindByName("Admin"); err != nil {
		t.Fatalf("Error occured %v", err)
	} else if role == nil {
		t.Fatal("Did not create Admin role")
	}

	if user, err := userRepo.FindByEmail("test@test.test"); err != nil {
		t.Fatalf("Error occured %v", err)
	} else if user == nil {
		t.Fatal("Did not create user")
	}
}
