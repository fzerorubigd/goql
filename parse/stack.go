package parse

import (
	"fmt"
	"sync"
)

// Stack is very silly implementation of stack, just for simplify things in where part
type Stack interface {
	Pop() (Item, error)
	Push(...Item)
	Peek() (Item, error)
}

type stack struct {
	s    []Item
	lock *sync.Mutex
}

func (s *stack) Pop() (Item, error) {
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

func (s *stack) Push(st ...Item) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.s = append(s.s, st...)
}

func (s *stack) Peek() (Item, error) {
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
		s:    make([]Item, 0, capacity),
		lock: &sync.Mutex{},
	}
}
