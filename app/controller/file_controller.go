package controller

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	. "github.com/netbrain/cloudfiler/app/repository"
	"strings"
	"time"
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
		Name:     name,
		Owner:    owner,
		Data:     data,
		Uploaded: time.Now(),
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

//TODO !!!NEEDS OPTIMIZATION!!!
func (c *FileController) FilesWhereUserHasAccess(user User) ([]File, error) {
	result := []File{}
	files, err := c.Files()
	if err != nil {
		return result, err
	}

FILES:
	for _, file := range files {
		if file.Owner.Equals(user) {
			result = append(result, file)
			continue FILES
		}

		for _, u := range file.Users {
			if u.Equals(user) {
				result = append(result, file)
				continue FILES
			}
		}

		for _, role := range file.Roles {
			if role.HasUser(user) {
				result = append(result, file)
				continue FILES
			}
		}
	}
	return result, nil
}

func (c *FileController) GrantUserAccessToFile(user User, file *File) error {
	if !c.UserHasAccess(user, *file) {
		file.Users = append(file.Users, user)
		return c.fileRepository.Store(file)
	}
	return nil
}

func (c *FileController) RevokeUserAccessToFile(user User, file *File) error {
	if c.UserHasAccess(user, *file) {
		for i, u := range file.Users {
			if user.Equals(u) {
				l := len(file.Users) - 1
				file.Users[i] = file.Users[l]
				file.Users = file.Users[:l]
				break
			}
		}
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

func (c *FileController) RevokeRoleAccessToFile(role Role, file *File) error {
	if c.RoleHasAccess(role, *file) {
		for i, r := range file.Roles {
			if role.Equals(r) {
				l := len(file.Roles) - 1
				file.Roles[i] = file.Roles[l]
				file.Roles = file.Roles[:l]
				break
			}
		}
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

func (c *FileController) AddTags(file *File, tags ...string) error {
TAGS:
	for _, tag := range tags {
		for _, existingTag := range file.Tags {
			if tag == existingTag {
				continue TAGS
			}
		}
		file.Tags = append(file.Tags, tag)
	}
	return c.fileRepository.Store(file)
}

func (c *FileController) RemoveTags(file *File, tags ...string) error {
	for _, tag := range tags {
		for i, t := range file.Tags {
			if tag == t {
				l := len(file.Tags) - 1
				file.Tags[i] = file.Tags[l]
				file.Tags = file.Tags[:l]
				break
			}
		}
	}
	return c.fileRepository.Store(file)
}

func (c *FileController) SetTags(file *File, tags ...string) error {
	file.Tags = []string{}
	return c.AddTags(file, tags...)
}

//TODO !!!NEEDS OPTIMIZATION!!!
//TODO add indexing of files upon app start and when adding/removing files and metadata
func (c *FileController) FileSearch(user User, query string) ([]File, error) {
	query = strings.ToLower(query)
	result := make([]File, 0)
	files, err := c.FilesWhereUserHasAccess(user)
	if err != nil {
		return files, err
	}

FILE:
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name), query) {
			result = append(result, file)
			continue
		}
		for _, tag := range file.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				result = append(result, file)
				continue FILE
			}
		}

	}
	return result, nil
}
