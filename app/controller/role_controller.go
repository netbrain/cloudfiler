package controller

import (
	"fmt"
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository"
	"strings"
)

type RoleController struct {
	roleRepository RoleRepository
	userRepository UserRepository
}

func NewRoleController(roleRepository RoleRepository, userRepository UserRepository) RoleController {
	c := RoleController{
		roleRepository: roleRepository,
		userRepository: userRepository,
	}
	return c
}

func (c RoleController) Roles() ([]Role, error) {
	return c.roleRepository.All()
}

func (c RoleController) Role(id int) (*Role, error) {
	return c.roleRepository.FindById(id)
}

func (c RoleController) RoleByName(name string) (*Role, error) {
	name = c.normalizeName(name)
	return c.roleRepository.FindByName(name)
}

func (c RoleController) Create(name string) (*Role, error) {
	name = c.normalizeName(name)
	if role, err := c.roleRepository.FindByName(name); role != nil {
		return nil, fmt.Errorf("Cannot create role, name already registered")
	} else if err != nil {
		return nil, err
	}

	r := &Role{Name: name}
	err := c.roleRepository.Store(r)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (c RoleController) Delete(id int) error {
	return c.roleRepository.Erase(id)
}

func (c RoleController) Update(id int, name string) error {
	name = c.normalizeName(name)

	if role, err := c.roleRepository.FindById(id); role != nil {
		role.Name = name
		return c.roleRepository.Store(role)
	} else if err != nil {
		return err
	}
	return fmt.Errorf("Role not found")
}

func (c RoleController) AddUser(role *Role, user *User) error {
	if !role.HasUser(*user) {
		role.Users = append(role.Users, *user)
		return c.roleRepository.Store(role)
	}
	return nil
}

func (c RoleController) RemoveUser(role *Role, user *User) error {
	for i, u := range role.Users {
		if user.Equals(u) {
			l := len(role.Users) - 1
			role.Users[i] = role.Users[l]
			role.Users = role.Users[:l]
			break
		}
	}
	return c.roleRepository.Store(role)
}

func (c RoleController) normalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.Title(name)
	return name
}
