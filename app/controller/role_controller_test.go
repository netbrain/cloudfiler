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
	roleController = NewRoleController(roleRepo)
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
	roleController.Create(name)
	role, _ := roleRepo.FindByName(name)
	if role == nil {
		t.Fatal("Role was not created!")
	}
}

func TestDeleteRole(t *testing.T) {
	initRoleControllerTest()

	name := "Testrole"
	roleController.Create(name)
	role, _ := roleRepo.FindByName(name)
	roleController.Delete(role.ID)

	all, _ := roleRepo.All()
	if l := len(all); l != 0 {
		t.Fatalf("No role should exist, but found %v (%#v)", l, all)
	}
}

func TestCreateTwoRolesWithIdenticalName(t *testing.T) {
	initRoleControllerTest()

	roleController.Create("Testrole")
	if err := roleController.Create("TestRole"); err == nil {
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
