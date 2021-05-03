package mockstore

import (
	"github.com/ashtanko/octo-server/store"
)

// Store ...
type Store struct {
	accountStore     *AccountStore
	transactionStore *TransactionStore
	jobStore         *JobStore
}

// Job mock store
func (s *Store) Job() store.JobStore {
	if s.jobStore != nil {
		return s.jobStore
	}

	s.jobStore = &JobStore{
		store: s,
	}

	return s.jobStore
}

func (s Store) Close() {
	// TODO
}

func (s *Store) Account() store.AccountStore {
	if s.accountStore != nil {
		return s.accountStore
	}

	s.accountStore = &AccountStore{
		store: s,
	}

	return s.accountStore
}

func (s *Store) Transaction() store.TransactionStore {
	if s.transactionStore != nil {
		return s.transactionStore
	}
	s.transactionStore = &TransactionStore{
		store: s,
	}

	return s.transactionStore
}

// New ...
func New() *Store {
	return &Store{}
}
