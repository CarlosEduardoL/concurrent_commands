package model

type stringStack struct {
	stack Stack
}

type StringStack interface {
	Push(string)
	Pop() (string, error)
	Size() int
}

func NewStringStack() StringStack {
	return &stringStack{NewStack()}
}

func (s *stringStack) Size() int {
	return s.stack.Size()
}

func (s *stringStack) Push(element string) {
	s.stack.Push(element)
}

func (s *stringStack) Pop() (string, error) {
	value, err := s.stack.Pop()
	return value.(string), err
}
