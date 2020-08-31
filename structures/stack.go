package structures

import (
	"errors"
	"sync"
)

type stack struct {
	mutLock sync.Mutex
	stack   []string
}

type Stack interface {
	Push(string)
	Pop() (string, error)
	Size() int
}

func NewStack() Stack {
	return &stack{sync.Mutex{}, make([]string, 0)}
}

func (s *stack) Size() int {
	return len(s.stack)
}

func (s *stack) Push(element string) {
	s.mutLock.Lock()
	defer s.mutLock.Unlock()

	s.stack = append(s.stack, element)
}

func (s *stack) Pop() (string, error) {
	s.mutLock.Lock()
	defer s.mutLock.Unlock()

	length := len(s.stack)
	if length == 0 {
		return "", errors.New("Empty stack")
	}

	element := s.stack[length-1]
	s.stack = s.stack[:length-1]
	return element, nil
}
