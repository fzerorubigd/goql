package parse

import (
	"fmt"
	"sync"
)

// Stack is very silly implementation of stack
type Stack interface {
	Pop() (item, error)
	Push(...item)
	Peek() (item, error)
}

type stack struct {
	s    []item
	lock *sync.Mutex
}

func (s *stack) Pop() (item, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if len(s.s) == 0 {
		return item{}, fmt.Errorf("stack is empty")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *stack) Push(st ...item) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.s = append(s.s, st...)
}

func (s *stack) Peek() (item, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if len(s.s) == 0 {
		return item{}, fmt.Errorf("stack is empty")
	}

	res := s.s[l-1]
	return res, nil
}

// NewStack create a new stack with initial capacity
func NewStack(capacity int) Stack {
	return &stack{
		s:    make([]item, 0, capacity),
		lock: &sync.Mutex{},
	}
}
