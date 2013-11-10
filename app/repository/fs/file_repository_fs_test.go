package fs

import (
	. "github.com/netbrain/cloudfiler/app/entity"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var fileRepo FileRepositoryFs

func cleanup() {
	os.RemoveAll(storagePath)
}

func initFileRepositoryFsTest() {
	userRepo = NewUserRepository()
	roleRepo = NewRoleRepository(userRepo)
	fileRepo = NewFileRepository(userRepo, roleRepo)
}

func createFile(id int, name string) *File {

	user := &User{
		ID:    1,
		Email: "test@test.test",
	}
	userRepo.Store(user)
	role := &Role{
		ID:    id,
		Name:  name,
		Users: []User{*user},
	}
	roleRepo.Store(role)

	file := &File{
		ID:       id,
		Name:     name,
		Uploaded: time.Now(),
		Data:     NewFileData(),
	}
	file.Owner = *user
	file.Roles = []Role{*role}
	file.Users = []User{*user}

	err := fileRepo.Store(file)

	if err != nil {
		panic(err)
	}

	return file
}

func TestGetFilePath(t *testing.T) {
	defer cleanup()
	initFileRepositoryFsTest()
	if p := fileRepo.getPath(0); p != storagePath+"/file/0" {
		t.Fatalf("Expected /tmp/file/0 but got %v", p)
	}
}

func TestStoreFile(t *testing.T) {
	defer cleanup()
	initFileRepositoryFsTest()
	bytes := []byte{1, 2, 3}
	file := createFile(1, "Test")
	file.Data.Write(bytes)

	_, err := ioutil.ReadFile(storagePath + "/file/1")

	if err != nil {
		t.Fatal(err)
	}

	_, err = ioutil.ReadFile(storagePath + "/filedata/1")

	if err != nil {
		t.Fatal(err)
	}

	if file.Data.Size() != 3 {
		t.Fatal("Expected 3 bytes")
	}
}

func TestUpdateFile(t *testing.T) {
	defer cleanup()
	initFileRepositoryFsTest()

	//create test file
	file := createFile(1, "Test")
	file.Description = "An updated description"
	fileRepo.Store(file)

	file, _ = fileRepo.FindById(file.ID)
	if file.Description != "An updated description" {
		t.Fatalf("Expected an updated description, instead it was: '%s'", file.Description)
	}
}

func TestEraseFile(t *testing.T) {
	defer cleanup()
	initFileRepositoryFsTest()
	file := createFile(1, "Test")

	err := fileRepo.Erase(file.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindFileById(t *testing.T) {
	defer cleanup()
	initFileRepositoryFsTest()
	file := createFile(1, "Test")

	file, err := fileRepo.FindById(file.ID)

	if err != nil {
		t.Fatal(err)
	}

	if file.ID != 1 ||
		file.Name != "Test" ||
		file.Owner.ID != 1 ||
		len(file.Roles) != 1 ||
		file.Roles[0].ID != 1 ||
		len(file.Users) != 1 ||
		file.Users[0].ID != 1 {

		t.Logf("%#v", file)
		t.Fatal("Inconsistent data")
	}
}

func TestFindAllFiles(t *testing.T) {
	defer cleanup()
	initFileRepositoryFsTest()
	createFile(1, "Test")

	files, err := fileRepo.All()

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatal("Expected 1 file")
	}
}
