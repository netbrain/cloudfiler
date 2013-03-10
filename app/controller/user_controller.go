package controller

import (
	"fmt"
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository"
	"strings"
)

type UserController struct {
	UserRepository UserRepository
}

func NewUserController(repository UserRepository) UserController {
	c := UserController{
		UserRepository: repository,
	}
	return c
}

func (c UserController) Users() ([]User, error) {
	return c.UserRepository.All()
}

func (c UserController) User(id int) (*User, error) {
	return c.UserRepository.FindById(id)
}

func (c UserController) UserByEmail(email string) (*User, error) {
	return c.UserRepository.FindByEmail(email)
}

func (c UserController) Create(email, password string) (*User, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	if user, err := c.UserRepository.FindByEmail(email); user != nil {
		return nil, fmt.Errorf("Cannot create user, email already registered")
	} else if err != nil {
		return nil, err
	}

	u := &User{Email: email}
	u.SetPassword(password)

	err := c.UserRepository.Store(u)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (c UserController) Delete(id int) error {
	return c.UserRepository.Erase(id)
}

func (c UserController) Update(id int, email, password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}
	if user, err := c.UserRepository.FindById(id); user != nil {
		user.SetPassword(password)
		user.Email = email
		return c.UserRepository.Store(user)
	} else if err != nil {
		return err
	}
	return fmt.Errorf("User not found")
}

func (c UserController) Count() (int, error) {
	return c.UserRepository.Count()
}

func validatePassword(password string) error {
	if len(password) < 3 {
		return fmt.Errorf("Password is required to be atleast 3 characters")
	}
	return nil
}
