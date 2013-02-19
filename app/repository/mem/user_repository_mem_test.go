package mem

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	"math/rand"
	"strconv"
	"testing"
)

var userRepo UserRepositoryMem

func initUserTest() {
	userRepo = NewUserRepositoryMem()
}

func getUser() *User {
	return &User{
		Email: strconv.Itoa(rand.Int()) + "@test.test",
	}
}

func TestUserStore(t *testing.T) {
	initUserTest()
	user := getUser()

	if err := userRepo.Store(user); err != nil {
		t.Fatal("Failed to store user")
	}
}

func TestUserStoreGeneratesRandomID(t *testing.T) {
	initUserTest()
	user := getUser()

	userRepo.Store(user)
	if user.ID == 0 {
		t.Fatal("Failed to generate id!")
	}
}

func TestUserErase(t *testing.T) {
	initUserTest()
	user := getUser()
	userRepo.Store(user)
	if err := userRepo.Erase(user.ID); err != nil {
		t.Fatalf("Error occured trying to erase user, %v", err)
	}
	if users, _ := userRepo.All(); len(users) != 0 {
		t.Fatal("User count should be 0")
	}
}

func TestUserGetAll(t *testing.T) {
	initUserTest()
	if users, _ := userRepo.All(); len(users) != 0 {
		t.Fatal("User count should be 0")
	}

	users := []*User{
		getUser(),
		getUser(),
		getUser(),
	}

	for _, user := range users {
		userRepo.Store(user)
	}

	if users, _ := userRepo.All(); len(users) != 3 {
		t.Fatalf("User count should be 3 but was %d", len(users))
	}

	storedUsers, _ := userRepo.All()
	for _, user1 := range users {
		found := false
		for _, user2 := range storedUsers {
			if user1.Equals(&user2) {
				found = true
			}
		}
		if !found {
			t.Fatalf("Did not find match in users. \nAll() returns: %v \n"+
				"while underlying map is: %v", storedUsers, userRepo.users)
		}
	}
}

func TestUserFindById(t *testing.T) {
	initUserTest()

	user := getUser()
	userRepo.Store(user)

	if user, err := userRepo.FindById(user.ID); err != nil {
		t.Fatalf("Error occured in finding user, %v", err)
	} else if !user.Equals(user) {
		t.Fatal("User doesn't equal stored user")
	}
}

func TestUserFindByEmail(t *testing.T) {
	initUserTest()

	user := getUser()
	userRepo.Store(user)

	if user, err := userRepo.FindByEmail(user.Email); err != nil {
		t.Fatalf("Error occured in finding user, %v", err)
	} else if !user.Equals(user) {
		t.Fatal("User doesn't equal stored user")
	}
}
