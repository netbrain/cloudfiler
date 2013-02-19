package mem

import (
	. "github.com/netbrain/cloudfiler/app/entity"
)

type UserRepositoryMem struct {
	users map[int]User
}

func NewUserRepository() UserRepositoryMem {
	return UserRepositoryMem{
		users: make(map[int]User),
	}
}

func (r UserRepositoryMem) Store(user *User) error {
	if user.ID == 0 {
		user.ID = generateID()
	}
	r.users[user.ID] = *user
	return nil
}

func (r UserRepositoryMem) Erase(id int) error {
	delete(r.users, id)
	return nil
}

func (r UserRepositoryMem) All() ([]User, error) {
	users := []User{}
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r UserRepositoryMem) FindById(id int) (*User, error) {
	if user, ok := r.users[id]; ok {
		return &user, nil
	}
	return nil, nil
}

func (r UserRepositoryMem) FindByEmail(email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}
