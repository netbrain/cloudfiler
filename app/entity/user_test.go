package entity

import (
	"testing"
)

func TestUserEquals(t *testing.T) {
	u1 := &User{ID: 1}
	u2 := &User{ID: 1}
	if !u1.Equals(u2) {
		t.Fatal("u1 is not equal to u2!")
	}
}

func TestUserNotEquals(t *testing.T) {
	u1 := &User{ID: 1}
	u2 := &User{ID: 2}
	if u1.Equals(u2) {
		t.Fatal("u1 is equal to u2!")
	}
}

func TestSetPasswordAndCheckHash(t *testing.T) {
	u := new(User)
	password := "test"
	if err := u.SetPassword(password); err != nil {
		t.Fatalf("Error occured when setting password: %v", err)
	}
	if !u.PasswordEquals(password) {
		t.Fatal("Password check failed")
	}
}
