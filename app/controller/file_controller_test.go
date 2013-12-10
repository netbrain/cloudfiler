package controller

import (
	. "github.com/netbrain/cloudfiler/app/repository/mem"
	"reflect"
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

func TestAddTags(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	file, _ := fileController.Create("filename.txt", *user, new(FileDataMem))
	file2, _ := fileController.Create("filename2.txt", *user, new(FileDataMem))

	fileController.AddTags(file, "TestTag", "TestTag2")
	fileController.AddTags(file2, "TestTag")

	if len(file.Tags) != 2 {
		t.Fatal("Expected 2 tags")
	}

	if len(fileController.tags) != 2 {
		t.Fatal("Expected two tags in tagmap")
	}

	if !reflect.DeepEqual(
		file.Tags,
		[]string{"TestTag", "TestTag2"},
	) {
		t.Fatal("Expected TestTag")
	}

	if fileController.tags["TestTag"] != 2 {
		t.Fatal("Expected 2 as tag count on TestTag")
	}
}

func TestRemoveTags(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	file, _ := fileController.Create("filename.txt", *user, new(FileDataMem))
	file2, _ := fileController.Create("filename2.txt", *user, new(FileDataMem))

	fileController.AddTags(file, "TestTag", "TestTag2")
	fileController.AddTags(file2, "TestTag")

	fileController.RemoveTags(file, "TestTag2")

	if len(file.Tags) != 1 {
		t.Fatal("Expected 1 tag")
	}

	if len(fileController.tags) != 1 {
		t.Fatal("Expected one tag in tagmap")
	}

	if !reflect.DeepEqual(
		file.Tags,
		[]string{"TestTag"},
	) {
		t.Fatal("Expected TestTag")
	}

	if fileController.tags["TestTag"] != 2 {
		t.Fatal("Expected 2 as tag count on TestTag")
	}

	if _, ok := fileController.tags["TestTag2"]; ok {
		t.Fatal("Did not expect the presence of TestTag2")
	}
}

func TestFileSearchName(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	fileController.Create("filename.txt", *user, new(FileDataMem))
	fileController.Create("filename2.txt", *user, new(FileDataMem))

	results, _ := fileController.FileSearch(*user, "filename")
	if len(results) != 2 {
		t.Fatal("expected two results for query filename")
	}
	results, _ = fileController.FileSearch(*user, "filename2")
	if len(results) != 1 {
		t.Fatal("expected one results for query: filename2")
	}
}

func TestFileSearchTag(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	file, _ := fileController.Create("filename.txt", *user, new(FileDataMem))
	file2, _ := fileController.Create("filename2.txt", *user, new(FileDataMem))

	fileController.AddTags(file, "TestTag", "TestTag2")
	fileController.AddTags(file2, "TestTag")

	results, _ := fileController.FileSearch(*user, "testtag")
	if len(results) != 2 {
		t.Fatal("expected two results for query: testtag")
	}
	results, _ = fileController.FileSearch(*user, "TestTag2")
	if len(results) != 1 {
		t.Fatal("expected one results for query: TestTag2")
	}
}

func TestFileSearchDescription(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	file, _ := fileController.Create("filename.txt", *user, new(FileDataMem))
	file2, _ := fileController.Create("filename2.txt", *user, new(FileDataMem))
	file.Description = "Description"
	file2.Description = "Description2"
	fileController.Update(file)
	fileController.Update(file2)

	results, _ := fileController.FileSearch(*user, "Description")
	if len(results) != 2 {
		t.Fatal("expected two results for query: Description")
	}
	results, _ = fileController.FileSearch(*user, "Description2")
	if len(results) != 1 {
		t.Fatal("expected one results for query: Description2")
	}
}

func TestFileSearchLongName(t *testing.T) {
	initFileControllerTest()
	user, _ := userController.Create("test@test.test", "testpasswd")
	fileController.Create("cookie recipie.txt", *user, new(FileDataMem))

	results, _ := fileController.FileSearch(*user, "cookie recipie")
	if len(results) != 1 {
		t.Fatal("expected one result")
	}
}
