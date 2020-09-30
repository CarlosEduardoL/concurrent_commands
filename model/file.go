package model

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type file struct {
	base string
	uri  string
	size int64
}

type File interface {
	Size() int64
	Remove() error
	Name() string
	String() string
	Copy(newBase string) error
}

func NewFile(base, uri string, size int64) File {
	return &file{base, uri, size}
}

func (f *file) Remove() error {
	return os.Remove(filepath.Join(f.base, f.uri))
}

func (f *file) Size() int64 {
	return f.size
}

func (f *file) Name() string {
	return filepath.Base(f.uri)
}

func (f *file) String() string {
	return f.uri
}

func (f *file) Copy(newBase string) error {
	sourceFile := filepath.Join(f.base, f.uri)
	destinationFile := filepath.Join(newBase, f.uri)
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	stat, err := os.Stat(filepath.Join(f.base, f.uri))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destinationFile, input, stat.Mode())
	if err != nil {
		return err
	}

	return nil
}
