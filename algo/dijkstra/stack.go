package dijkstra

import "errors"

type Stack struct {
	stack []int
}

func NewStack() *Stack {
	return &Stack{
		stack: make([]int, 0),
	}
}

func (s *Stack) Push(v int) {
	s.stack = append(s.stack, v)
}

func (s *Stack) Pop() (int, error) {
	l := len(s.stack)
	if l == 0 {
		return 0, errors.New("Stack is empty")
	}
	rv := 0
	s.stack, rv = s.stack[:l-1], s.stack[l-1]
	return rv, nil
}

func (s *Stack) AsSlice() []int {
	cp := make([]int, len(s.stack))
	copy(cp, s.stack)
	return cp
}
