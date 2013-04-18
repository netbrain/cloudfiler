package fs

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	"io/ioutil"
	"testing"
)

var userRepo UserRepositoryFs

func initUserRepositoryFsTest() {
	userRepo = NewUserRepository()
}

func createUser(id int, email string) *User {
	user := &User{
		ID:    id,
		Email: email,
	}

	err := userRepo.Store(user)

	if err != nil {
		panic(err)
	}

	return user
}

func TestGetUserPath(t *testing.T) {
	defer cleanup()
	initUserRepositoryFsTest()
	if p := userRepo.getPath(0); p != storagePath+"/user/0" {
		t.Fatalf("Expected %s/user/0 but got %v", storagePath, p)
	}
}

func TestStoreUser(t *testing.T) {
	defer cleanup()
	initUserRepositoryFsTest()
	createUser(1, "test@test.test")

	_, err := ioutil.ReadFile(storagePath + "/user/1")

	if err != nil {
		t.Fatal(err)
	}
}

func TestEraseUser(t *testing.T) {
	defer cleanup()
	initUserRepositoryFsTest()
	user := createUser(1, "test@test.test")

	err := userRepo.Erase(user.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindUserById(t *testing.T) {
	defer cleanup()
	initUserRepositoryFsTest()
	user := createUser(1, "test@test.test")

	user, err := userRepo.FindById(user.ID)

	if err != nil {
		t.Fatal(err)
	}

	if !(user.ID == 1 && user.Email == "test@test.test") {
		t.Logf("%#v", user)
		t.Fatal("Inconsistent data")
	}
}

func TestFindAllUsers(t *testing.T) {
	defer cleanup()
	initUserRepositoryFsTest()
	createUser(1, "test@test.test")

	users, err := userRepo.All()

	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Fatal("Expected 1 user")
	}
}

func TestCountUsers(t *testing.T) {
	defer cleanup()
	initUserRepositoryFsTest()
	createUser(1, "test@test.test")

	count, err := userRepo.Count()
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatal("Expected 1")
	}
}
