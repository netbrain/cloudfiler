package fs

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	"io/ioutil"
)

type UserRepositoryFs struct{}

func NewUserRepository() UserRepositoryFs {
	return UserRepositoryFs{}
}

func (r UserRepositoryFs) Store(user *User) error {
	if user.ID == 0 {
		user.ID = generateID()
	}

	path := r.getPath(user.ID)
	err := ioutil.WriteFile(path, serialize(user), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepositoryFs) All() ([]User, error) {
	users := make([]User, 0)
	fileList, err := ioutil.ReadDir(r.getPath(""))
	if err != nil {
		return nil, err
	}

	for _, fi := range fileList {
		if !fi.IsDir() {
			path := r.getPath(fi.Name())
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}

			user := &User{}
			unserialize(b, user)
			users = append(users, *user)
		}
	}
	return users, nil
}

func (r UserRepositoryFs) FindById(id int) (*User, error) {
	b, err := ioutil.ReadFile(r.getPath(id))
	if err != nil {
		return nil, err
	}

	user := &User{}
	unserialize(b, user)

	return user, nil
}

func (r UserRepositoryFs) FindByEmail(email string) (*User, error) {
	users, err := r.All()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (r UserRepositoryFs) Count() (int, error) {
	fileList, err := ioutil.ReadDir(r.getPath(""))
	if err != nil {
		return 0, err
	}

	return len(fileList), nil
}

func (r UserRepositoryFs) getPath(id interface{}) string {
	return getPath("user", id)
}
