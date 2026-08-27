package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type identity interface {
	Username() string
	Role() string
}

type operatorIdentity struct {
	op *Operator
}

func (i operatorIdentity) Username() string {
	if i.op == nil {
		return ""
	}
	return i.op.Username
}

func (i operatorIdentity) Role() string {
	if i.op == nil {
		return ""
	}
	return i.op.Role
}

func principalOf(op *Operator) identity {
	if op == nil {
		return nil
	}
	return operatorIdentity{op: op}
}

type sessionKey struct{}

func loginHandler(ops *OperatorStore, sessions *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid login JSON"})
			return
		}
		op, err := ops.Verify(body.Username, body.Password)
		if err != nil {
			opsJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}
		ident := principalOf(op)
		if ident == nil {
			opsJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}
		session := sessions.Issue(ident.Username(), ident.Role())
		opsJSON(w, http.StatusOK, map[string]any{"token": session.Token, "username": ident.Username(), "role": ident.Role(), "expires_at": session.ExpiresAt})
	}
}

func requireAuth(sessions *SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimSpace(r.Header.Get("X-Session-Token"))
			if token == "" {
				opsJSON(w, http.StatusUnauthorized, map[string]string{"error": "token required"})
				return
			}
			session, ok := sessions.Verify(token)
			if !ok {
				opsJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
				return
			}
			ctx := context.WithValue(r.Context(), sessionKey{}, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func meHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(sessionKey{}).(Session)
		if !ok {
			opsJSON(w, http.StatusUnauthorized, map[string]string{"error": "not authenticated"})
			return
		}
		opsJSON(w, http.StatusOK, map[string]any{"username": session.Username, "role": session.Role})
	}
}
