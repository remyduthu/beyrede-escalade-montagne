package beyredeescalademontagne

import (
	"net/http"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DBDir string
}

type server struct {
	db *badger.DB
}

func New(cfg Config) (*server, error) {
	db, err := openDB(cfg.DBDir)
	if err != nil {
		return nil, err
	}

	result := &server{db}
	addHandlers(result)

	return result, nil
}

func addHandlers(s *server) {
	http.Handle("/state", s.handleState())
	http.Handle("/schedules", s.handleSchedules())
}

func openDB(dir string) (*badger.DB, error) {
	opts := badger.DefaultOptions(dir)
	opts.Logger = logrus.StandardLogger()

	return badger.Open(opts)
}
