package mem

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	"reflect"
	"testing"
)

var roleRepo RoleRepositoryMem

func initRoleTest() {
	roleRepo = NewRoleRepositoryMem()
}

func TestRoleStore(t *testing.T) {
	initRoleTest()
	role := &Role{
		ID:    generateID(),
		Name:  "TestRole",
		Users: []User{*getUser()},
	}

	roleRepo.Store(role)
	if v, ok := roleRepo.roles[role.ID]; !ok {
		t.Fatal("Role not stored")
	} else {
		if !reflect.DeepEqual(role, &v) {
			t.Fatal("Stored role vs local role has mismatch")
		}
	}
}

func TestRoleErase(t *testing.T) {
	initRoleTest()
	role := &Role{
		ID: generateID(),
	}
	roleRepo.Store(role)
	if len(roleRepo.roles) != 1 {
		t.Fatal("Role was not stored")
	}
	roleRepo.Erase(role.ID)
	if len(roleRepo.roles) != 0 {
		t.Fatal("Role was not erased")
	}
}

func TestRoleAll(t *testing.T) {
	initRoleTest()
	role := &Role{}
	roleRepo.Store(role)
	roles, _ := roleRepo.All()

	if len(roles) != 1 {
		t.Fatal("Expected 1 result")
	}

	if !reflect.DeepEqual(role, &roles[0]) {
		t.Fatal("Expected equality")
	}
}

func TestRoleFindById(t *testing.T) {
	initRoleTest()
	role := &Role{
		ID: generateID(),
	}

	roleRepo.Store(role)
	r, _ := roleRepo.FindById(role.ID)

	if !reflect.DeepEqual(role, r) {
		t.Fatal("Expected equality")
	}
}

func TestRoleFindByName(t *testing.T) {
	initRoleTest()
	role := &Role{
		Name: "TestRole",
	}

	roleRepo.Store(role)
	r, _ := roleRepo.FindByName("TestRole")

	if !reflect.DeepEqual(role, r) {
		t.Fatal("Expected equality")
	}
}
