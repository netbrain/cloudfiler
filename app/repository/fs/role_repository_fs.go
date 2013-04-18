package fs

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	"io/ioutil"
	"os"
	"strconv"
)

type RoleRepositoryFs struct{}

type RoleFs struct {
	ID    int
	Name  string
	Users []int
}

func NewRoleRepository() RoleRepositoryFs {
	return RoleRepositoryFs{}
}

func (r RoleRepositoryFs) Store(role *Role) error {
	if role.ID == 0 {
		role.ID = generateID()
	}

	data := RoleFs{
		ID:    role.ID,
		Name:  role.Name,
		Users: make([]int, 0),
	}

	for _, user := range role.Users {
		data.Users = append(data.Users, user.ID)
	}

	path := r.getPath(role.ID)
	err := ioutil.WriteFile(path, serialize(data), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (r RoleRepositoryFs) Erase(id int) error {
	path := r.getPath(id)
	return os.Remove(path)
}

func (r RoleRepositoryFs) All() ([]Role, error) {
	roles := make([]Role, 0)
	fileList, err := ioutil.ReadDir(r.getPath(""))
	if err != nil {
		return nil, err
	}

	for _, fi := range fileList {
		if !fi.IsDir() {
			id, err := strconv.Atoi(fi.Name())
			if err != nil {
				return nil, err
			}
			role, err := r.FindById(id)
			if err != nil {
				return nil, err
			}
			roles = append(roles, *role)
		}
	}
	return roles, nil
}

func (r RoleRepositoryFs) FindById(id int) (*Role, error) {
	b, err := ioutil.ReadFile(r.getPath(id))
	if err != nil {
		return nil, err
	}

	rolefs := &RoleFs{}
	unserialize(b, rolefs)

	role := &Role{
		ID:   rolefs.ID,
		Name: rolefs.Name,
		//TODO add users
	}

	return role, nil
}

func (r RoleRepositoryFs) FindByName(name string) (*Role, error) {
	roles, err := r.All()
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.Name == name {
			return &role, nil
		}
	}
	return nil, nil
}

func (r RoleRepositoryFs) getPath(id interface{}) string {
	return getPath("role", id)
}
