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
	fileRepo = NewFileRepository()
}

func createFile(id int, name string) *File {
	file := &File{
		ID:       id,
		Name:     name,
		Uploaded: time.Now(),
		Data:     NewFileData(),
	}

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

	if !(file.ID == 1 && file.Name == "Test") {
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
