package mem

import (
	. "github.com/netbrain/cloudfiler/app/entity"
)

type RoleRepositoryMem struct {
	roles map[int]Role
}

func NewRoleRepository() RoleRepositoryMem {
	return RoleRepositoryMem{
		roles: make(map[int]Role),
	}
}

func (r RoleRepositoryMem) Store(role *Role) error {
	if role.ID == 0 {
		role.ID = generateID()
	}
	r.roles[role.ID] = *role
	return nil
}

func (r RoleRepositoryMem) Erase(id int) error {
	delete(r.roles, id)
	return nil
}

func (r RoleRepositoryMem) All() ([]Role, error) {
	roles := []Role{}
	for _, role := range r.roles {
		roles = append(roles, role)
	}
	return roles, nil
}

func (r RoleRepositoryMem) FindById(id int) (*Role, error) {
	if role, ok := r.roles[id]; ok {
		return &role, nil
	}
	return nil, nil
}

func (r RoleRepositoryMem) FindByName(name string) (*Role, error) {
	for _, role := range r.roles {
		if role.Name == name {
			return &role, nil
		}
	}
	return nil, nil
}
