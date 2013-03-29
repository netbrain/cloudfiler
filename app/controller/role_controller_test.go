package controller

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	"testing"
)

var roleRepo RoleRepositoryMem
var roleController RoleController

func initRoleControllerTest() {
	roleRepo = NewRoleRepository()
	userRepo = NewUserRepository()
	roleController = NewRoleController(roleRepo, userRepo)
	userController = NewUserController(userRepo)
}

func TestRoles(t *testing.T) {
	initRoleControllerTest()
	roles, _ := roleController.Roles()
	if len(roles) != 0 {
		t.Fatal("Expected zero roles")
	}

	roleController.Create("testrole")

	roles, _ = roleController.Roles()
	if len(roles) != 1 {
		t.Fatal("Expected one role")
	}
}

func TestCreateRole(t *testing.T) {
	initRoleControllerTest()

	name := "Testrole"
	role, _ := roleController.Create(name)
	if role == nil {
		t.Fatal("Role was not created!")
	}
}

func TestDeleteRole(t *testing.T) {
	initRoleControllerTest()

	name := "Testrole"
	role, _ := roleController.Create(name)
	roleController.Delete(role.ID)

	all, _ := roleRepo.All()
	if l := len(all); l != 0 {
		t.Fatalf("No role should exist, but found %v (%#v)", l, all)
	}
}

func TestCreateTwoRolesWithIdenticalName(t *testing.T) {
	initRoleControllerTest()

	roleController.Create("Testrole")
	if _, err := roleController.Create("TestRole"); err == nil {
		t.Fatal("Illegal creation of two roles with identical name! this should have failed!")
	}
}

func TestGetRoleByID(t *testing.T) {
	initRoleControllerTest()
	role := &Role{
		Name: "Testrole",
	}
	roleRepo.Store(role)
	role2, err := roleController.Role(role.ID)

	if err != nil {
		t.Fatalf("Error occured trying to get role by id, %v", err)
	}

	if role2 == nil {
		t.Fatal("Role not found! recieved nil")
	}
}

func TestGetRoleByIDWhereNoneExist(t *testing.T) {
	initRoleControllerTest()
	role, err := roleController.Role(0)

	if err != nil {
		t.Fatalf("Error occured trying to get role by id, %v", err)
	}

	if role != nil {
		t.Fatal("Role found! when it should not!")
	}
}

func TestAddUserToRole(t *testing.T) {
	initRoleControllerTest()
	role, err := roleController.Create("testrole")

	if err != nil {
		t.Fatal("Error occured when trying to create role")
	}

	user, err := userController.Create("test@test.test", "testpasswd")

	if err != nil {
		t.Fatal("Error occured when trying to create user")
	}

	if err := roleController.AddUser(role, user); err != nil {
		t.Fatalf("Got error when trying to add user to role %v", err)
	}

	role, err = roleController.Role(role.ID)
	if err != nil {
		t.Fatal("Error occured when trying to fetch  updated role")
	}

	if len(role.Users) != 1 {
		t.Fatalf("Expected 1 user added")
	}

	if !user.Equals(role.Users[0]) {
		t.Fatal("Expected equality")
	}
}

func TestRoleHasUser(t *testing.T) {
	initRoleControllerTest()
	role, _ := roleController.Create("testrole")
	user, _ := userController.Create("test@test.test", "testpasswd")

	if role.HasUser(*user) {
		t.Fatal("Fresh role should not have a user")
	}

	roleController.AddUser(role, user)

	if !role.HasUser(*user) {
		t.Fatal("user not found in role!")
	}

}

func TestRemoveUserFromRole(t *testing.T) {
	initRoleControllerTest()
	role, _ := roleController.Create("testrole")
	user, _ := userController.Create("test@test.test", "testpasswd")
	roleController.AddUser(role, user)

	if len(role.Users) != 1 {
		t.Fatalf("Expected 1 user added")
	}

	if err := roleController.RemoveUser(role, user); err != nil {
		t.Fatal(err)
	}

	if len(role.Users) != 0 {
		t.Fatalf("Expected zero users")
	}
}
