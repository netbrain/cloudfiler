package fs

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	"io/ioutil"
	"testing"
)

var roleRepo RoleRepositoryFs

func initRoleRepositoryFsTest() {
	roleRepo = NewRoleRepository()
}

func createRole(id int, name string) *Role {
	role := &Role{
		ID:   id,
		Name: name,
	}

	err := roleRepo.Store(role)

	if err != nil {
		panic(err)
	}

	return role
}

func TestGetRolePath(t *testing.T) {
	defer cleanup()
	initRoleRepositoryFsTest()
	if p := roleRepo.getPath(0); p != storagePath+"/role/0" {
		t.Fatalf("Expected %s/role/0 but got %v", storagePath, p)
	}
}

func TestStoreRole(t *testing.T) {
	defer cleanup()
	initRoleRepositoryFsTest()
	createRole(1, "Test")

	_, err := ioutil.ReadFile(storagePath + "/role/1")

	if err != nil {
		t.Fatal(err)
	}
}

func TestEraseRole(t *testing.T) {
	defer cleanup()
	initRoleRepositoryFsTest()
	role := createRole(1, "Test")

	err := roleRepo.Erase(role.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindRoleById(t *testing.T) {
	defer cleanup()
	initRoleRepositoryFsTest()
	role := createRole(1, "Test")

	role, err := roleRepo.FindById(role.ID)

	if err != nil {
		t.Fatal(err)
	}

	if !(role.ID == 1 && role.Name == "Test") {
		t.Logf("%#v", role)
		t.Fatal("Inconsistent data")
	}
}

func TestFindAllRoles(t *testing.T) {
	defer cleanup()
	initRoleRepositoryFsTest()
	createRole(1, "Test")

	roles, err := roleRepo.All()

	if err != nil {
		t.Fatal(err)
	}

	if len(roles) != 1 {
		t.Fatal("Expected 1 role")
	}
}
