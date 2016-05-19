package entity

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    string
	Password []byte
}

func (u *User) Equals(other interface{}) bool {
	switch o := other.(type) {
	case User:
		return u.ID == o.ID
	case *User:
		return u.ID == o.ID
	}
	return false
}

func (u *User) SetPassword(password string) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = encrypted
	return nil
}

func (u *User) PasswordEquals(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password)) == nil
}

func (u User) String() string {
	return u.Email
}
