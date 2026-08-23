package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var sessionSequence uint64

type Session struct {
	Token     string
	Username  string
	Role      string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]Session
	clock    OpsClock
	ttl      time.Duration
}

func newSessionStore() *SessionStore {
	return &SessionStore{clock: newOpsClock(), ttl: time.Hour}
}

func (s *SessionStore) Issue(username, role string) Session {
	now := s.clock.Now().UTC()
	token := fmt.Sprintf("tok-%06d", atomic.AddUint64(&sessionSequence, 1))
	session := Session{Token: token, Username: username, Role: role, IssuedAt: now, ExpiresAt: now.Add(s.ttl)}
	s.mu.Lock()
	s.sessions[token] = session
	s.mu.Unlock()
	return session
}

func (s *SessionStore) Verify(token string) (Session, bool) {
	s.mu.RLock()
	session, ok := s.sessions[token]
	s.mu.RUnlock()
	if !ok {
		return Session{}, false
	}
	if s.clock.Now().UTC().After(session.ExpiresAt) {
		return Session{}, false
	}
	return session, true
}

func (s *SessionStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}
