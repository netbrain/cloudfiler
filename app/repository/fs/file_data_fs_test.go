package fs

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestTempFile(t *testing.T) {
	f, _ := ioutil.TempFile("", "")
	defer os.Remove(f.Name()) //cleanup

	if _, err := os.Stat(f.Name()); os.IsNotExist(err) {
		t.Fatal("File doesn't exist")
	}

	f.Close()

	if _, err := os.Stat(f.Name()); os.IsNotExist(err) {
		t.Fatal("File doesn't exist, it must have been autoremoved after Close() which shouldnt happen")
	}

}

func TestNewFileData(t *testing.T) {
	fd := NewFileData()
	defer os.Remove(fd.file.Name()) //cleanup

	if fd.Size() != 0 {
		t.Fatal("Expected new zero sized temporary file")
	}
}

func TestWriteToNewFileData(t *testing.T) {
	fd := NewFileData()
	defer os.Remove(fd.file.Name()) //cleanup

	written, err := fd.Write([]byte{1, 2, 3})

	if err != nil {
		t.Fatalf("Recieved error %v", err)
	}

	if written != 3 {
		t.Fatal("Expected 3 bytes written")
	}
}
