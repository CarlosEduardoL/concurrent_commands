package model

import (
	"errors"
	"sync"
)

type stack struct {
	mutLock sync.Mutex
	head    *node
	size    int
}

type Stack interface {
	Push(interface{})
	Pop() (interface{}, error)
	Size() int
}

func NewStack() Stack {
	return &stack{
		mutLock: sync.Mutex{},
		head:    nil,
		size:    0}
}

func (s *stack) Size() int {
	return s.size
}

func (s *stack) Push(element interface{}) {
	s.mutLock.Lock()
	defer s.mutLock.Unlock()
	s.head = &node{next: s.head, value: element}
	s.size++
}

func (s *stack) Pop() (interface{}, error) {
	s.mutLock.Lock()
	defer s.mutLock.Unlock()

	if s.size == 0 {
		return nil, errors.New("Empty stack")
	}

	s.size--

	element := s.head.value
	s.head = s.head.next
	return element, nil
}
