package mem

import (
	. "github.com/netbrain/cloudfiler/app/entity"
)

type FileRepositoryMem struct {
	files map[int]File
}

func NewFileRepository() FileRepositoryMem {
	return FileRepositoryMem{
		files: make(map[int]File),
	}
}

func (r FileRepositoryMem) Store(file *File) error {
	if file.ID == 0 {
		file.ID = generateID()
	}
	r.files[file.ID] = *file
	return nil
}

func (r FileRepositoryMem) Erase(id int) error {
	delete(r.files, id)
	return nil
}

func (r FileRepositoryMem) All() ([]File, error) {
	files := []File{}
	for _, file := range r.files {
		files = append(files, file)
	}
	return files, nil
}

func (r FileRepositoryMem) FindById(id int) (*File, error) {
	if file, ok := r.files[id]; ok {
		return &file, nil
	}
	return nil, nil
}
