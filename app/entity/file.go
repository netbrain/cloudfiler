package entity

import (
	"time"
)

type FileData interface {
	Close() error
	Read([]byte) (int, error)
	Size() int64
	Write([]byte) (int, error)
	Seek(int64, int) (int64, error)
}

type File struct {
	ID          int
	Name        string
	Owner       User
	Tags        []string
	Users       []User
	Roles       []Role
	Data        FileData
	Uploaded    time.Time
	Description string
}

func (f *File) Equals(other interface{}) bool {
	switch o := other.(type) {
	case File:
		return f.ID == o.ID
	case *File:
		return f.ID == o.ID
	}
	return false
}
