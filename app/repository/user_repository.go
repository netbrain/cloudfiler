package repository

import (
	. "github.com/netbrain/cloudfiler/app/entity"
)

type UserRepository interface {
	Store(user *User) error
	Erase(id int) error
	All() ([]User, error)
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	Count() (int, error)
}
