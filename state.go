package beyredeescalademontagne

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type state string

const (
	busyState   state = "busy"
	closedState state = "closed"
	openedState state = "opened"
)

func (s state) validate() error {
	switch s {
	case busyState, closedState, openedState:
		return nil
	default:
		return errors.New("invalid state")
	}
}

type timedState struct {
	State state     `json:"state"`
	Until time.Time `json:"until"`
}

func (s *timedState) dbKey() []byte { return []byte("state") }

func (s *timedState) validate() error {
	if err := s.State.validate(); err != nil {
		return err
	}

	if s.Until.Before(time.Now().Add(time.Hour)) {
		return errors.New("invalid until")
	}

	return nil
}

func (s *server) handleState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleGetState(w)
		case http.MethodPatch:
			s.handlePatchState(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (s *server) handleGetState(w http.ResponseWriter) {
	state, err := getState(s.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	data, err := json.Marshal(state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getState(db *badger.DB) (*timedState, error) {
	result := &timedState{}

	return result, db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(result.dbKey())
		if err != nil {
			return err
		}

		return item.Value(func(value []byte) error {
			return json.Unmarshal(value, item)
		})
	})
}

func (s *server) handlePatchState(w http.ResponseWriter, r *http.Request) {
	if err := foo(s.db, r.Header.Get("Authorization")); err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	params := &timedState{}
	if err := json.NewDecoder(r.Body).Decode(params); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := updateState(s.db, params); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// TODO(remyduthu): Rename function.
func foo(db *badger.DB, authz string) error {
	strings.TrimPrefix(authz, "Bearer ")

	authzParts := strings.Split(authz, ".")
	if len(authzParts) != 2 { //nolint:gomnd
		return errors.New("")
	}

	username, password := authzParts[0], authzParts[1]

	user, err := getUser(db, username)
	if err != nil {
		return errors.New("")
	}

	if password != user.Password {
		return errors.New("")
	}

	return nil
}

func updateState(db *badger.DB, state *timedState) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(state)
		if err != nil {
			return err
		}

		return txn.Set(state.dbKey(), data)
	})
}
