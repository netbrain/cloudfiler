package controller

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	"testing"
)

var userRepo UserRepositoryMem
var userController UserController

func initUserControllerTest() {
	userRepo = NewUserRepository()
	userController = NewUserController(userRepo)
}

func TestGetAllUsersWhenNoneCreated(t *testing.T) {
	initUserControllerTest()

	if result, _ := userController.Users(); len(result) != 0 {
		t.Fatal("There should be no users!")
	}
}

func TestCreateUser(t *testing.T) {
	initUserControllerTest()

	email := "create@user.test"
	userController.Create(email, "password")
	user, _ := userRepo.FindByEmail(email)
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
		userController.Create(email, "test")
	}

	for _, email := range emails {
		user, _ := userRepo.FindByEmail(email)
		if user == nil {
			t.Fatal("User was not created!")
		}
	}
}

func TestCreateUserWithInvalidPassword(t *testing.T) {
	initUserControllerTest()

	email := "create@user.test"
	if _, err := userController.Create(email, "12"); err == nil {
		t.Fatal("Password should not be valid")
	}
}

func TestCreateTwoUsersWithIdenticalEmail(t *testing.T) {
	initUserControllerTest()

	userController.Create("dup@dup.test", "pwd")
	if _, err := userController.Create("dup@dup.test", "pwd"); err == nil {
		t.Fatal("Illegal creation of two users with identical email! this should have failed!")
	}
}

func TestGetUserByID(t *testing.T) {
	initUserControllerTest()
	user1 := &User{
		Email: "test@test.test",
	}
	userRepo.Store(user1)
	user2, err := userController.User(user1.ID)

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
	user, err := userController.User(0)

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
	userRepo.Store(user1)
	user2, err := userController.UserByEmail(user1.Email)

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
	user, err := userController.UserByEmail("test@test.test")

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
	userRepo.Store(user)

	email := "new@email.test"
	pwd := "anotherpwd"
	if err := userController.Update(user.ID, email, pwd); err != nil {
		t.Fatalf("update returned in error %v", err)
	}

	user, _ = userRepo.FindById(user.ID)
	if user.Email != email {
		t.Fatal("email not updated")
	}

	if !user.PasswordEquals(pwd) {
		t.Fatal("password not updated")
	}
}

func TestUpdateUserWhereUserDoNotExist(t *testing.T) {
	initUserControllerTest()
	if err := userController.Update(0, "Some@email.test", "passwd"); err == nil {
		t.Fatal("update did not return in error when it should")
	}
}

func TestUpdateUserWithInvalidPassword(t *testing.T) {
	initUserControllerTest()
	user := &User{
		Email: "test@test.test",
	}
	userRepo.Store(user)

	pwd := "--"
	if err := userController.Update(user.ID, user.Email, pwd); err == nil {
		t.Fatal("update should return in error")
	}

	if user.PasswordEquals(pwd) {
		t.Fatal("password should not have been updated")
	}
}
