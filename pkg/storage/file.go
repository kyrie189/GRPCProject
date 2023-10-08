package storage

import (
	"bytes"
)

type File struct {
	Name   string
	buffer *bytes.Buffer
}

func NewFile(name string) *File {
	return &File{
		Name:   name,
		buffer: &bytes.Buffer{},
	}
}

func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)

	return err
}
