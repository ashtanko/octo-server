package mockstore

import (
	"github.com/ashtanko/octo-server/model"
)

type TransactionStore struct {
	store *Store
}

// Fetch list of transaction
func (t *TransactionStore) Fetch(count int) ([]model.Transaction, error) {
	var list []model.Transaction
	return list, nil
}

// Update the transaction record to CANCELLED
func (t *TransactionStore) Update(id int) error {
	return nil
}

// Save the transaction
func (t *TransactionStore) Save(request *model.IncomingRequest, status string) error {
	return nil
}
