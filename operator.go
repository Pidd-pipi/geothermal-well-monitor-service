package main

import (
	"fmt"
	"sync"
)

type Operator struct {
	Username string
	Name     string
	Role     string
	Password string
}

type OperatorStore struct {
	mu     sync.RWMutex
	byName map[string]Operator
}

func newOperatorStore(seed []Operator) *OperatorStore {
	s := &OperatorStore{byName: map[string]Operator{}}
	for _, op := range seed {
		s.byName[op.Username] = op
	}
	return s
}

func (s *OperatorStore) Lookup(username string) (*Operator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	op, ok := s.byName[username]
	if !ok {
		return nil, nil
	}
	copy := op
	return &copy, nil
}

func (s *OperatorStore) Verify(username, password string) (*Operator, error) {
	op, err := s.Lookup(username)
	if err != nil {
		return nil, err
	}
	if op == nil {
		return op, nil
	}
	if op.Password != password {
		return nil, fmt.Errorf("%w: bad credentials", ErrOpsPolicy)
	}
	return op, nil
}

func (s *OperatorStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.byName)
}
