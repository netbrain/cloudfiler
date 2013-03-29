package controller

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository"
)

type FileController struct {
	fileRepository FileRepository
}

func NewFileController(fileRepository FileRepository) FileController {
	return FileController{
		fileRepository: fileRepository,
	}
}

func (c *FileController) Create(name string, owner User, data FileData) (*File, error) {
	file := &File{
		Name:  name,
		Owner: owner,
		Data:  data,
	}

	if err := c.fileRepository.Store(file); err != nil {
		return nil, err
	}

	return file, nil
}

func (c *FileController) Erase(id int) error {
	return c.fileRepository.Erase(id)
}

func (c *FileController) Files() ([]File, error) {
	return c.fileRepository.All()
}

func (c *FileController) File(id int) (*File, error) {
	return c.fileRepository.FindById(id)
}

func (c *FileController) FilesWhereUserHasAccess(user User) ([]File, error) {
	return c.fileRepository.AllWhereUserHasAccess(user)
}

func (c *FileController) GrantUserAccessToFile(user User, file *File) error {
	if !c.UserHasAccess(user, *file) {
		file.Users = append(file.Users, user)
		return c.fileRepository.Store(file)
	}
	return nil
}

func (c *FileController) GrantRoleAccessToFile(role Role, file *File) error {
	if !c.RoleHasAccess(role, *file) {
		file.Roles = append(file.Roles, role)
		return c.fileRepository.Store(file)
	}
	return nil
}

func (c *FileController) UserHasAccess(user User, file File) bool {
	if file.Owner.Equals(user) {
		return true
	}

	for _, r := range file.Roles {
		if r.HasUser(user) {
			return true
		}
	}

	for _, u := range file.Users {
		if user.Equals(u) {
			return true
		}
	}

	return false
}

func (c *FileController) RoleHasAccess(role Role, file File) bool {
	for _, r := range file.Roles {
		if role.Equals(r) {
			return true
		}
	}
	return false
}
