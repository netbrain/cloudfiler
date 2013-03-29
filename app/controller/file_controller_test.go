package controller

import (
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	"testing"
)

var fileController FileController
var fileRepo FileRepositoryMem

func initFileControllerTest() {
	userRepo = NewUserRepository()
	roleRepo = NewRoleRepository()
	fileRepo = NewFileRepository()
	userController = NewUserController(userRepo)
	roleController = NewRoleController(roleRepo, userRepo)
	fileController = NewFileController(fileRepo)

}

func TestCreateFile(t *testing.T) {
	initFileControllerTest()
	fileData := new(FileDataMem)
	user, _ := userController.Create("test@test.test", "testpasswd")
	file, _ := fileController.Create("filename.txt", *user, fileData)
	files, _ := fileController.Files()
	if len(files) != 1 {
		t.Fatal("File was not created")
	}

	if !file.Equals(files[0]) {
		t.Fatalf("Expected equal file")
	}
}

func TestRetrieveOnlyFilesOwnedByUser(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	otherUser, _ := userController.Create("other@test.test", "testpasswd")
	fileController.Create("filename.txt", *otherUser, new(FileDataMem))
	file, _ := fileController.Create("filename2.txt", *user, new(FileDataMem))
	files, _ := fileController.FilesWhereUserHasAccess(*user)
	if len(files) != 1 {
		t.Fatal("Should be one file in list")
	}

	if !file.Equals(files[0]) {
		t.Fatalf("Expected equal file")
	}
}

func TestGrantUserAccessToFile(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	otherUser, _ := userController.Create("other@test.test", "testpasswd")
	file, _ := fileController.Create("filename.txt", *user, new(FileDataMem))
	fileController.GrantUserAccessToFile(*otherUser, file)
	if !fileController.UserHasAccess(*otherUser, *file) {
		t.Fatal("User doesn't have access")
	}
}

func TestGrantRoleAccessToFile(t *testing.T) {
	initFileControllerTest()
	otherUser, _ := userController.Create("other@test.test", "testpasswd")
	user, _ := userController.Create("test@test.test", "testpasswd")
	role, _ := roleController.Create("TestRole")
	file, _ := fileController.Create("filename.txt", *otherUser, new(FileDataMem))

	roleController.AddUser(role, user)
	fileController.GrantRoleAccessToFile(*role, file)

	if !fileController.RoleHasAccess(*role, *file) {
		t.Fatal("Role should have access")
	}

	if !fileController.UserHasAccess(*user, *file) {
		t.Fatal("User should have access")
	}
}

func TestEraseFile(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	file, _ := fileController.Create("filename.txt", *user, new(FileDataMem))
	files, _ := fileController.Files()
	if len(files) != 1 {
		t.Fatal("File was not created")
	}

	fileController.Erase(file.ID)
	files, _ = fileController.Files()
	if len(files) != 0 {
		t.Fatal("File was not erased")
	}
}
