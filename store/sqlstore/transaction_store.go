package sqlstore

import (
	"github.com/ashtanko/octo-server/model"
	"github.com/sirupsen/logrus"
)

type TransactionStore struct {
	store *Store
}

// Fetch list of transaction
func (t *TransactionStore) Fetch(count int) ([]model.Transaction, error) {
	var list []model.Transaction
	var tr model.Transaction
	sqlStatement := `
		SELECT id, amount FROM transactions WHERE status_key = 'DONE' ORDER BY timestamptz DESC LIMIT $1 FOR UPDATE;
	`

	if err := t.store.Db.Ping(); err != nil {
		return list, err
	}
	rows, err := t.store.Db.Query(sqlStatement, count)
	if err != nil {
		return list, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}()
	for rows.Next() {
		err := rows.Scan(&tr.Id, &tr.Amount)
		if err != nil {
			logrus.Fatal(err)
		}
		list = append(list, tr)
	}
	return list, nil
}

// Update the transaction record to CANCELLED
func (t *TransactionStore) Update(id int) error {
	sqlStatement := `
		UPDATE transactions SET status_key = 'CANCELLED' WHERE id = $1;
	`

	if err := t.store.Db.Ping(); err != nil {
		return err
	}

	_, err := t.store.Db.Exec(sqlStatement, id)

	if err != nil {
		return err
	}

	return nil
}

// Save the transaction
func (t *TransactionStore) Save(request *model.IncomingRequest, status string) error {

	sqlStatement := `
		INSERT INTO transactions (type_key, status_key, amount, transactionId) 
		VALUES ($1, $2, $3, $4) 
		RETURNING *
	`

	if err := t.store.Db.Ping(); err != nil {
		return err
	}

	return t.store.Db.QueryRow(
		sqlStatement,
		request.Source,
		status,
		request.Amount,
		request.TransactionID,
	).Scan(
		&request.Source,
		&request.State,
		&request.Amount,
		&request.TransactionID,
	)
}
