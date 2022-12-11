package beyredeescalademontagne

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *user) dbKey() []byte {
	return []byte("user-" + u.Username)
}

func (s *server) UpdateUser(username, password string) error {
	return updateUser(s.db, &user{username, password})
}

func getUser(db *badger.DB, username string) (*user, error) {
	result := &user{Username: username}

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

func updateUser(db *badger.DB, user *user) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return txn.Set(user.dbKey(), data)
	})
}
