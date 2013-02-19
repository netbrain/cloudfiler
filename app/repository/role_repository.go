package repository

import (
	. "github.com/netbrain/cloudfiler/app/entity"
)

type RoleRepository interface {
	Store(role *Role) error
	Erase(id int) error
	All() ([]Role, error)
	FindById(id int) (*Role, error)
	FindByName(name string) (*Role, error)
}
