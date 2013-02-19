package controller

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	"testing"
)

var repo UserRepositoryMem
var controller UserController

func initUserControllerTest() {
	repo = NewUserRepository()
	controller = NewUserController(repo)
}

func TestGetAllUsersWhenNoneCreated(t *testing.T) {
	initUserControllerTest()

	if result, _ := controller.Users(); len(result) != 0 {
		t.Fatal("There should be no users!")
	}
}

func TestCreateUser(t *testing.T) {
	initUserControllerTest()

	email := "create@user.test"
	controller.Create(email, "password")
	user, _ := repo.FindByEmail(email)
	if user == nil {
		t.Fatal("User was not created!")
	}
}

func TestCreateSeveralUsers(t *testing.T) {
	initUserControllerTest()

	emails := []string{
		"test1@test.test",
		"test2@test.test",
		"test3@test.test",
	}

	for _, email := range emails {
		controller.Create(email, "test")
	}

	for _, email := range emails {
		user, _ := repo.FindByEmail(email)
		if user == nil {
			t.Fatal("User was not created!")
		}
	}
}

func TestCreateUserWithInvalidPassword(t *testing.T) {
	initUserControllerTest()

	email := "create@user.test"
	if err := controller.Create(email, "12"); err == nil {
		t.Fatal("Password should not be valid")
	}
}

func TestDeleteUser(t *testing.T) {
	initUserControllerTest()

	email := "delete@user.test"
	controller.Create(email, "password")
	user, _ := repo.FindByEmail(email)
	controller.Delete(user.ID)

	all, _ := repo.All()
	if l := len(all); l != 0 {
		t.Fatalf("No user should exist, but found %v (%#v)", l, all)
	}
}

func TestCreateTwoUsersWithIdenticalEmail(t *testing.T) {
	initUserControllerTest()

	controller.Create("dup@dup.test", "pwd")
	if err := controller.Create("dup@dup.test", "pwd"); err == nil {
		t.Fatal("Illegal creation of two users with identical email! this should have failed!")
	}
}

func TestGetUserByID(t *testing.T) {
	initUserControllerTest()
	user1 := &User{
		Email: "test@test.test",
	}
	repo.Store(user1)
	user2, err := controller.User(user1.ID)

	if err != nil {
		t.Fatalf("Error occured trying to get user by id, %v", err)
	}

	if user2 == nil {
		t.Fatal("User not found! recieved nil")
	}

	if !user1.Equals(user2) {
		t.Fatal("user1 doesnt equal user2")
	}
}

func TestGetUserByIDWhereNoneExist(t *testing.T) {
	initUserControllerTest()
	user, err := controller.User(0)

	if err != nil {
		t.Fatalf("Error occured trying to get user by id, %v", err)
	}

	if user != nil {
		t.Fatal("User found! when it should not!")
	}
}

func TestGetUserByEmail(t *testing.T) {
	initUserControllerTest()
	user1 := &User{
		Email: "test@test.test",
	}
	repo.Store(user1)
	user2, err := controller.UserByEmail(user1.Email)

	if err != nil {
		t.Fatalf("Error occured trying to get user by id, %v", err)
	}

	if user2 == nil {
		t.Fatal("User not found! recieved nil")
	}

	if !user1.Equals(user2) {
		t.Fatal("user1 doesnt equal user2")
	}
}

func TestGetUserByEmailWhereNoneExist(t *testing.T) {
	initUserControllerTest()
	user, err := controller.UserByEmail("test@test.test")

	if err != nil {
		t.Fatalf("Error occured trying to get user by id, %v", err)
	}

	if user != nil {
		t.Fatal("User found! when it should not!")
	}
}

func TestUpdateUser(t *testing.T) {
	initUserControllerTest()
	user := &User{
		Email: "test@test.test",
	}
	repo.Store(user)

	email := "new@email.test"
	pwd := "anotherpwd"
	if err := controller.Update(user.ID, email, pwd); err != nil {
		t.Fatalf("update returned in error %v", err)
	}

	user, _ = repo.FindById(user.ID)
	if user.Email != email {
		t.Fatal("email not updated")
	}

	if !user.PasswordEquals(pwd) {
		t.Fatal("password not updated")
	}
}

func TestUpdateUserWhereUserDoNotExist(t *testing.T) {
	initUserControllerTest()
	if err := controller.Update(0, "Some@email.test", "passwd"); err == nil {
		t.Fatal("update did not return in error when it should")
	}
}

func TestUpdateUserWithInvalidPassword(t *testing.T) {
	initUserControllerTest()
	user := &User{
		Email: "test@test.test",
	}
	repo.Store(user)

	pwd := "--"
	if err := controller.Update(user.ID, user.Email, pwd); err == nil {
		t.Fatal("update should return in error")
	}

	if user.PasswordEquals(pwd) {
		t.Fatal("password should not have been updated")
	}
}
