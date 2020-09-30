package model

import (
	"errors"
	"sync"
)

type fileNode struct {
	next  *fileNode
	value File
}

type fileStack struct {
	mutLock sync.Mutex
	head    *fileNode
	size    int
}

type FileStack interface {
	Push(File)
	Pop() (File, error)
	Size() int
}

func NewFileStack() FileStack {
	return &fileStack{
		mutLock: sync.Mutex{},
		head:    nil,
		size:    0}
}

func (s *fileStack) Size() int {
	return s.size
}

func (s *fileStack) Push(element File) {
	s.mutLock.Lock()
	defer s.mutLock.Unlock()
	s.head = &fileNode{s.head, element}

}

func (s *fileStack) Pop() (File, error) {
	s.mutLock.Lock()
	defer s.mutLock.Unlock()

	if s.size == 0 {
		return nil, errors.New("Empty stack")
	}

	element := s.head.value
	s.head = s.head.next
	return element, nil
}
