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

type ByUploaded []File

func (f *File) Equals(other interface{}) bool {
	switch o := other.(type) {
	case File:
		return f.ID == o.ID
	case *File:
		return f.ID == o.ID
	}
	return false
}

func (f *File) FormattedUploaded() string {
	return f.Uploaded.Format(time.RFC822)
}

func (u ByUploaded) Len() int {
	return len(u)
}

func (u ByUploaded) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u ByUploaded) Less(i, j int) bool {
	return u[j].Uploaded.Before(u[i].Uploaded)
}
