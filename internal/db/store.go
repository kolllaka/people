package db

import (
	"database/sql"

	"github.com/KoLLlaka/sobes/internal/logger"
)

type Store interface {
	NewPeopleStore() PeopleStore
}

type store struct {
	db     *sql.DB
	logger logger.MyLogger
}

// create a new Store interface
func NewStore(db *sql.DB, logger logger.MyLogger) Store {
	return &store{
		db:     db,
		logger: logger,
	}
}

func (s *store) NewPeopleStore() PeopleStore {
	return &peopleStore{
		db:     s.db,
		logger: s.logger,
	}
}
