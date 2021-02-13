package mockstore

import (
	"github.com/ashtanko/octo-server/model"
)

// AccountStore ...
type AccountStore struct {
	store *Store
}

// GetBalance Get the account balance
func (r *AccountStore) GetBalance(id string) (float64, error) {
	u := &model.Account{}
	return u.AccountBalance, nil
}

// UpdateBalance Increase the account balance
func (r *AccountStore) UpdateBalance(id string, balance float64) error {
	return nil
}

// SetBalance account balance by id
func (r *AccountStore) SetBalance(id string, balance float64) error {
	return nil
}

// Create Account with balance = 0
func (r *AccountStore) Create(u *model.Account) error {
	return nil
}

// Find Return account by id
func (r *AccountStore) Find(id string) (*model.Account, error) {
	u := &model.Account{}
	return u, nil
}

// Delete account by id
func (r *AccountStore) Delete(id string) error {
	return nil
}
