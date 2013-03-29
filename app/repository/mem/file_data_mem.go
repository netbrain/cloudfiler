package mem

import (
	"bytes"
	"errors"
	"os"
)

type FileDataMem struct {
	pos    int64
	data   []byte
	length int64
}

func (f *FileDataMem) Close() error {
	f.pos = 0
	return nil
}

func (f *FileDataMem) Read(b []byte) (int, error) {
	bLen := int64(len(b))
	read := 0
	for p := int64(0); p < bLen; p++ {
		pos := p + f.pos
		if pos >= f.Size() {
			break
		}

		b[p] = f.data[pos]
		read++
	}
	f.pos += int64(read)
	return read, nil
}

func (f *FileDataMem) Size() int64 {
	return f.length
}

func (f *FileDataMem) Write(b []byte) (int, error) {
	bLen := len(b)
	f.data = bytes.Join([][]byte{
		f.data[:f.pos],
		b,
		f.data[f.pos:],
	}, []byte{})

	f.length += int64(bLen)
	f.pos += int64(bLen)

	return bLen, nil
}

func (f *FileDataMem) Seek(offset int64, whence int) (int64, error) {
	var pos int64
	switch whence {
	case os.SEEK_SET:
		pos = offset
	case os.SEEK_CUR:
		pos += offset
	case os.SEEK_END:
		pos = f.length + offset
	}

	if pos > f.Size() || pos < 0 {
		return f.pos, errors.New("seek is out of range")
	}

	f.pos = pos
	return f.pos, nil
}
