package fs

import (
	"fmt"
	"os"
)

type FileDataFs struct {
	file     *os.File
	isClosed bool
}

func NewFileData() *FileDataFs {
	f := &FileDataFs{
		file: getTempFile(),
	}

	f.Close()

	return f
}

func (f *FileDataFs) Close() error {
	f.isClosed = true
	return f.file.Close()
}

func (f *FileDataFs) Read(b []byte) (int, error) {
	f.reopen()
	return f.file.Read(b)
}

func (f *FileDataFs) Size() int64 {
	f.reopen()
	fi, err := f.file.Stat()
	if err != nil {
		panic(err)
	}
	return fi.Size()
}

func (f *FileDataFs) Write(b []byte) (int, error) {
	f.reopen()
	written, err := f.file.Write(b)
	if err != nil {
		f.file.Sync()
	}
	return written, err
}

func (f *FileDataFs) Seek(offset int64, whence int) (ret int64, err error) {
	f.reopen()
	ret, err = f.file.Seek(offset, whence)
	fmt.Printf("%v", err)
	return
}

func (f *FileDataFs) Rename(to string) error {
	var err error
	if err = os.Rename(f.file.Name(), to); err != nil {
		return err
	}
	f.file, err = os.OpenFile(to, os.O_RDWR|os.O_EXCL, 0600)
	f.isClosed = false
	return err
}

func (f *FileDataFs) reopen() {
	if f.isClosed {
		var err error
		f.file, err = os.OpenFile(f.file.Name(), os.O_RDWR|os.O_EXCL, 0600)
		if err != nil {
			panic(err)
		}
		f.isClosed = false
	}
}
