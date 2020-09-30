package model

import (
	"os"
	"path/filepath"
)

type folder struct {
	base string
	uri  string
}

type Folder interface {
	Name() string
	Remove() error
	Create() error
	ReBase(newBase string)
}

func NewFolder(base, uri string) Folder {
	return &folder{base, uri}
}

func (f *folder) Name() string {
	return filepath.Base(f.uri)
}

func (f *folder) Remove() error {
	return os.Remove(filepath.Join(f.base, f.uri))
}

func (f *folder) Create() error {
	stat, err := os.Stat(filepath.Join(f.base, f.uri))
	if err != nil {
		return err
	}
	return os.Mkdir(filepath.Join(f.base, f.uri), stat.Mode())
}

func (f *folder) ReBase(newBase string) {
	f.base = newBase
}
