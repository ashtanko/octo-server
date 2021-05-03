package sqlstore

import (
	"database/sql"
	"github.com/ashtanko/octo-server/model"
	"github.com/ashtanko/octo-server/store"
)

// AccountStore ...
type AccountStore struct {
	Store *Store
}

// GetBalance Get the account balance
func (r *AccountStore) GetBalance(id string) (float64, error) {
	sqlStatement := `
		SELECT account_balance FROM account WHERE id = $1;
	`

	u := &model.Account{}
	if err := r.Store.Db.QueryRow(
		sqlStatement,
		id,
	).Scan(
		&u.AccountBalance,
	); err != nil {
		if err == sql.ErrNoRows {
			return 0, store.ErrRecordNotFound
		}
		return 0, err
	}

	return u.AccountBalance, nil
}

// UpdateBalance Increase the account balance
func (r *AccountStore) UpdateBalance(id string, balance float64) error {

	sqlStatement := `
		UPDATE account SET account_balance = account_balance + $1 WHERE id = $2;
	`

	if err := r.Store.Db.Ping(); err != nil {
		return err
	}

	_, err := r.Store.Db.Exec(sqlStatement, balance, id)

	if err != nil {
		return err
	}

	return nil
}

// SetBalance account balance by id
func (r *AccountStore) SetBalance(id string, balance float64) error {

	sqlStatement := `
		UPDATE account SET account_balance = $1 WHERE id = $2;
	`

	if err := r.Store.Db.Ping(); err != nil {
		return err
	}

	_, err := r.Store.Db.Exec(sqlStatement, balance, id)

	if err != nil {
		return err
	}

	return nil
}

// Create Account with balance = 0
func (r *AccountStore) Create(u *model.Account) error {

	sqlStatement := `
		INSERT INTO account (account_balance) VALUES ($1) RETURNING id;
	`

	if err := r.Store.Db.Ping(); err != nil {
		return err
	}

	return r.Store.Db.QueryRow(
		sqlStatement,
		u.AccountBalance,
	).Scan(&u.ID)
}

// Find Return account by id
func (r *AccountStore) Find(id string) (*model.Account, error) {
	sqlStatement := `
		SELECT id, account_balance FROM account WHERE id = $1;
	`
	u := &model.Account{}
	if err := r.Store.Db.QueryRow(
		sqlStatement,
		id,
	).Scan(
		&u.ID,
		&u.AccountBalance,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

// Delete account by id
func (r *AccountStore) Delete(id string) error {

	sqlStatement := `
		DELETE FROM account
		WHERE id = $1;
	`

	_, err := r.Store.Db.Exec(sqlStatement, id)

	if err != nil {
		return store.ErrRecordNotFound
	}

	return nil
}
