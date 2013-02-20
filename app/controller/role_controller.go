package controller

import (
	"fmt"
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository"
	"strings"
)

type RoleController struct {
	RoleRepository RoleRepository
}

func NewRoleController(repository RoleRepository) RoleController {
	c := RoleController{
		RoleRepository: repository,
	}
	return c
}

func (c RoleController) Roles() ([]Role, error) {
	return c.RoleRepository.All()
}

func (c RoleController) Role(id int) (*Role, error) {
	return c.RoleRepository.FindById(id)
}

func (c RoleController) RoleByName(name string) (*Role, error) {
	name = c.normalizeName(name)
	return c.RoleRepository.FindByName(name)
}

func (c RoleController) Create(name string) error {
	name = c.normalizeName(name)
	if role, err := c.RoleRepository.FindByName(name); role != nil {
		return fmt.Errorf("Cannot create role, name already registered")
	} else if err != nil {
		return err
	}

	r := &Role{Name: name}

	return c.RoleRepository.Store(r)
}

func (c RoleController) Delete(id int) error {
	return c.RoleRepository.Erase(id)
}

func (c RoleController) Update(id int, name string) error {
	name = c.normalizeName(name)

	if role, err := c.RoleRepository.FindById(id); role != nil {
		role.Name = name
		return c.RoleRepository.Store(role)
	} else if err != nil {
		return err
	}
	return fmt.Errorf("Role not found")
}

func (c RoleController) normalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.Title(name)
	return name
}
