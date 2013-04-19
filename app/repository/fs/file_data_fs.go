package fs

import (
	"os"
)

type FileDataFs struct {
	file *os.File
}

func NewFileData() *FileDataFs {
	file := &FileDataFs{
		file: getTempFile(),
	}

	return file
}

func (f *FileDataFs) Close() error {
	return f.file.Close()
}

func (f *FileDataFs) Read(b []byte) (int, error) {
	return f.file.Read(b)
}

func (f *FileDataFs) Size() int64 {
	fi, err := f.file.Stat()
	if err != nil {
		panic(err)
	}
	return fi.Size()
}

func (f *FileDataFs) Write(b []byte) (int, error) {
	written, err := f.file.Write(b)
	if err != nil {
		f.file.Sync()
	}
	return written, err
}

func (f *FileDataFs) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}
