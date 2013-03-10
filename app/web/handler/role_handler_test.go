package handler

import (
	. "github.com/netbrain/cloudfiler/app/controller"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	. "github.com/netbrain/cloudfiler/app/web"
	. "github.com/netbrain/cloudfiler/app/web/testing"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

var roleRepo RoleRepositoryMem
var roleController RoleController
var roleHandler RoleHandler

func initRoleHandlerTest() {
	roleRepo = NewRoleRepository()
	userRepo = NewUserRepository()
	roleController = NewRoleController(roleRepo, userRepo)
	userController = NewUserController(userRepo)
	roleHandler = NewRoleHandler(roleController, userController)
}

func TestGetRoleList(t *testing.T) {
	initRoleHandlerTest()
	ctx, _ := CreateReqContext("GET", "/role/list", nil)
	result := roleHandler.List(ctx)
	if result == nil {
		t.Fatalf("response returned %v", result)
	}

}
func TestGetCreateRolePage(t *testing.T) {
	initRoleHandlerTest()
	ctx, _ := CreateReqContext("GET", "/role/create", nil)
	result := roleHandler.Create(ctx)
	if result != nil {
		t.Fatal("Expected nil")
	}
}

func TestCreateRole(t *testing.T) {
	initRoleHandlerTest()
	ctx, _ := CreateReqContext("POST", "/role/create", map[string][]string{
		"name": []string{"testrole"},
	})
	roleHandler.Create(ctx)

	if ctx.HasValidationErrors() {
		t.Log(ctx.ValidationErrors)
		t.Fatal("Did not expect any errors")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

	result, _ := roleRepo.All()
	if len(result) != 1 {
		t.Fatal("Expected role created")
	}
}

func TestCreateRoleWithNoParameters(t *testing.T) {
	initRoleHandlerTest()
	ctx, _ := CreateReqContext("POST", "/role/create", nil)

	roleHandler.Create(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}
}

func TestCreateRoleWhereRoleAlreadyExist(t *testing.T) {
	initRoleHandlerTest()
	roleController.Create("testrole")

	ctx, _ := CreateReqContext("POST", "/role/create", map[string][]string{
		"name": []string{"testrole"},
	})

	result := roleHandler.Create(ctx)

	switch result.(type) {
	case error:
		//Expects a application error
	default:
		t.Fatal("Unexpected return type")
	}
}

func TestRetrieveNonExistingRole(t *testing.T) {
	initRoleHandlerTest()

	ctx, _ := CreateReqContext("GET", "/role/retrieve", map[string][]string{
		"id": []string{"123"},
	})

	result := roleHandler.Retrieve(ctx)

	if err, ok := result.(*AppError); ok {
		if err.Status() != http.StatusNotFound {
			t.Fatal("Expected 404")
		}
	} else {
		t.Fatal("Expected app error")
	}
}

func TestRetrieveExistingRole(t *testing.T) {
	initRoleHandlerTest()

	role, _ := roleController.Create("testrole")

	ctx, _ := CreateReqContext("GET", "/role/retrieve?id="+strconv.Itoa(role.ID), nil)

	result := roleHandler.Retrieve(ctx)

	if !role.Equals(result) {
		t.Fatalf("Expected equality between\n%#v\nand\n%#v", role, result)
	}
}

func TestUpdateRole(t *testing.T) {
	initRoleHandlerTest()
	role, _ := roleController.Create("testrole")

	ctx, _ := CreateReqContext("POST", "/role/update", map[string][]string{
		"id":   []string{strconv.Itoa(role.ID)},
		"name": []string{"newtestrole"},
	})
	roleHandler.Update(ctx)

	if ctx.HasValidationErrors() {
		t.Log(ctx.ValidationErrors)
		t.Fatal("Did not expect any errors")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

	role, _ = roleRepo.FindById(role.ID)
	if role.Name != "Newtestrole" {
		t.Fatalf("Expected updated name, instead got: %#v", role)
	}
}

func TestUpdateRoleWithNoParameters(t *testing.T) {
	initRoleHandlerTest()
	role, _ := roleController.Create("testrole")
	ctx, _ := CreateReqContext("POST", "/role/update?id="+strconv.Itoa(role.ID), nil)

	roleHandler.Update(ctx)

	if !ctx.HasValidationErrors() {
		t.Fatal("Expected validation errors")
	}
}

func TestDeleteNonExistingRole(t *testing.T) {
	initRoleHandlerTest()
	ctx, _ := CreateReqContext("GET", "/role/delete?id=123", nil)

	result := roleHandler.Delete(ctx)

	if result != nil {
		t.Fatal("expected nil")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

}

func TestDeleteRole(t *testing.T) {
	initRoleHandlerTest()
	role, _ := roleController.Create("testrole")
	ctx, _ := CreateReqContext("GET", "/role/delete?id="+strconv.Itoa(role.ID), nil)

	result := roleHandler.Delete(ctx)

	if result != nil {
		t.Fatal("expected nil")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

	allRoles, _ := roleRepo.All()
	if len(allRoles) != 0 {
		t.Fatal("Expected zero roles")
	}
}

func TestRemoveUserFromRole(t *testing.T) {
	initRoleHandlerTest()
	role, _ := roleController.Create("testrole")
	user, _ := userController.Create("test@test.test", "testpasswd")

	if err := roleController.AddUser(role, user); err != nil {
		t.Fatal("Error occured trying to add user to role")
	}

	if len(role.Users) != 1 {
		t.Fatal("Expected one user")
	}

	ctx, _ := CreateReqContext("GET",
		"/role/users/remove?id="+
			strconv.Itoa(role.ID)+
			"&uid="+
			strconv.Itoa(user.ID),
		nil)

	result := roleHandler.RemoveUser(ctx)

	if result != nil {
		t.Fatal("expected nil")
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

	role, _ = roleController.Role(role.ID)
	if len(role.Users) != 0 {
		t.Fatal("Expected zero users")
	}
}

func TestGetAddUserToRole(t *testing.T) {
	initRoleHandlerTest()
	role, _ := roleController.Create("testrole")
	userController.Create("test@test.test", "testpasswd")
	ctx, _ := CreateReqContext("GET", "/role/users/add?id="+strconv.Itoa(role.ID), nil)

	result := roleHandler.AddUser(ctx)

	v := reflect.ValueOf(result)
	rolev := v.FieldByName("Role")
	usersv := v.FieldByName("Users")

	if !role.Equals(rolev.Interface()) {
		t.Fatal("Expected role")
	}

	if usersv.Len() != 1 {
		t.Fatal("Expected one user")
	}

}

func TestAddUserToRole(t *testing.T) {
	initRoleHandlerTest()
	role, _ := roleController.Create("testrole")
	user, _ := userController.Create("test@test.test", "testpasswd")

	ctx, _ := CreateReqContext("POST", "/role/users/add",
		map[string][]string{
			"id":  []string{strconv.Itoa(role.ID)},
			"uid": []string{strconv.Itoa(user.ID)},
		})

	result := roleHandler.AddUser(ctx)

	if result != nil {
		t.Fatalf("expected nil but got", result)
	}

	if !ctx.IsRedirected() {
		t.Fatal("Expected redirection")
	}

	role, _ = roleController.Role(role.ID)
	l := len(role.Users)
	if l != 1 {
		t.Fatalf("Expected one user instead got %s", l)
	}
}

func TestRemoveLastUserFromAdminRole(t *testing.T) {
	initRoleHandlerTest()
	role, _ := roleController.Create("Admin")
	user, _ := userController.Create("test@test.test", "testpasswd")

	roleController.AddUser(role, user)

	ctx, _ := CreateReqContext("GET",
		"/role/users/remove?id="+
			strconv.Itoa(role.ID)+
			"&uid="+
			strconv.Itoa(user.ID),
		nil)

	roleHandler.RemoveUser(ctx)

	role, _ = roleController.Role(role.ID)
	if len(role.Users) != 1 {
		t.Fatal("Expected one user")
	}
}
