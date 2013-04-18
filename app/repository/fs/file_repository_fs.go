package fs

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type FileRepositoryFs struct{}

type FileFs struct {
	ID       int
	Name     string
	Owner    int
	Tags     []string
	Users    []int
	Roles    []int
	Uploaded time.Time
}

func NewFileRepository() FileRepositoryFs {
	return FileRepositoryFs{}
}

func (r FileRepositoryFs) Store(file *File) error {
	var err error
	if file.ID == 0 {
		file.ID = generateID()
	}

	data := FileFs{
		ID:       file.ID,
		Name:     file.Name,
		Owner:    file.Owner.ID,
		Tags:     file.Tags,
		Users:    make([]int, 0),
		Roles:    make([]int, 0),
		Uploaded: file.Uploaded,
	}

	for _, user := range file.Users {
		data.Users = append(data.Users, user.ID)
	}

	for _, role := range file.Roles {
		data.Roles = append(data.Roles, role.ID)
	}

	path := r.getPath(file.ID)
	err = ioutil.WriteFile(path, serialize(data), 0600)
	if err != nil {
		return err
	}

	fileData := file.Data.(*FileDataFs)
	oldFileDataPath := fileData.file.Name()
	newFileDataPath := getPath("filedata", file.ID)

	err = os.Rename(oldFileDataPath, newFileDataPath)
	if err != nil {
		return err
	}

	fileData.file, err = os.OpenFile(newFileDataPath, os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (r FileRepositoryFs) Erase(id int) error {
	path := r.getPath(id)
	return os.Remove(path)
}

func (r FileRepositoryFs) All() ([]File, error) {
	files := make([]File, 0)
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
			file, err := r.FindById(id)
			if err != nil {
				return nil, err
			}
			files = append(files, *file)
		}
	}
	return files, nil
}

func (r FileRepositoryFs) FindById(id int) (*File, error) {
	b, err := ioutil.ReadFile(r.getPath(id))
	if err != nil {
		return nil, err
	}

	filefs := &FileFs{}
	unserialize(b, filefs)

	osfile, err := os.Open(getPath("filedata", id))
	if err != nil {
		return nil, err
	}

	file := &File{
		ID:   filefs.ID,
		Name: filefs.Name,
		//Owner: filefs.Owner,
		Tags:  filefs.Tags,
		Users: make([]User, 0),
		Roles: make([]Role, 0),
		Data: &FileDataFs{
			file: osfile,
		},
		Uploaded: filefs.Uploaded,
	}

	return file, nil
}

func (r FileRepositoryFs) getPath(id interface{}) string {
	return getPath("file", id)
}
