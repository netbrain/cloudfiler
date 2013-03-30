package repository

import (
	. "github.com/netbrain/cloudfiler/app/entity"
)

type FileRepository interface {
	Store(*File) error
	Erase(int) error
	All() ([]File, error)
	FindById(int) (*File, error)
}
