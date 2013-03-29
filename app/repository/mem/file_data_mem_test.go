package mem

import (
	"bytes"
	"os"
	"testing"
)

func TestWriteData(t *testing.T) {
	b := []byte{1, 2, 3}
	fd := new(FileDataMem)
	written, _ := fd.Write(b)
	if written != len(b) {
		t.Fatalf("Expected %s bytes to be written", len(b))
	}
	if !bytes.Equal(fd.data, b) {
		t.Fatalf("Byte arrays don't match, expected %v but got %v", b, fd.data)
	}

	if fd.Size() != int64(len(b)) {
		t.Fatalf("Expected Size() to return %s, instead got %s", len(b), fd.Size())
	}
}

func TestReadDataWhenTruncated(t *testing.T) {
	fd := new(FileDataMem)
	b := make([]byte, 4096)
	read, _ := fd.Read(b)

	if read != 0 {
		t.Fatal("Expected zero to be read")
	}
}

func TestReadDataWhenHasData(t *testing.T) {
	fd := new(FileDataMem)
	w := []byte{1, 2, 3}
	fd.Write(w)
	fd.Close()

	b := make([]byte, 4096)
	read, _ := fd.Read(b)

	if read != 3 {
		t.Fatalf("Expected %v bytes to be read, instead read is %v", len(w), read)
	}

	if !bytes.Equal(w, b[:read]) {
		t.Fatalf("Byte arrays don't match, expected %v but got %v", w, b[:read])
	}
}

func TestCanSeek(t *testing.T) {
	fd := new(FileDataMem)
	w := []byte{1, 2, 3}
	fd.Write(w)
	fd.Close()

	_, err := fd.Seek(0, os.SEEK_END)
	if err != nil {
		t.Fatal("seeker can't seek")
	}
	_, err = fd.Seek(0, os.SEEK_SET)
	if err != nil {
		t.Fatal("seeker can't seek")
	}
}
