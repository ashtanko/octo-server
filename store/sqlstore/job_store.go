package sqlstore

import (
	"github.com/ashtanko/octo-server/model"
	"github.com/sirupsen/logrus"
	"time"
)

type JobStore struct {
	Store *Store
}

func (j *JobStore) CancelTransactionAndCorrectBalance(jobRunInterval time.Duration, accountId int) (err error) {
	tx, err := j.Store.Db.Begin()
	if err != nil {
		logrus.Error(err)
		return
	}

	lastJobActivity := &time.Time{}
	sqlStatement := `
		SELECT last_activity_at FROM jobs WHERE id = $1 FOR UPDATE;
	`

	// fetch last job activity timestamp
	err = tx.QueryRow(sqlStatement, accountId).Scan(lastJobActivity)
	if err != nil {
		return tx.Rollback()
	}

	// check last successful job time, if less than job interval, just return and rollback sql transaction
	nTimeAgo := time.Now().Add(-jobRunInterval)
	if nTimeAgo.Before(*lastJobActivity) {
		logrus.Info("Current job not finished yet")
		return tx.Rollback()
	}

	sqlStatement = `
		SELECT id, amount FROM transactions WHERE status_key = 'DONE' ORDER BY timestamptz DESC LIMIT $1 FOR UPDATE;
	`

	// fetch order to cancel
	rows, err := tx.Query(sqlStatement, 20)
	if err != nil {
		return tx.Rollback()
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}()
	// no rows, nothing to do
	if !rows.Next() {
		return tx.Rollback()
	}

	var list []model.Transaction
	for rows.Next() {
		var tr model.Transaction
		err := rows.Scan(&tr.Id, &tr.Amount)
		if err != nil {
			logrus.Fatal(err)
		}
		if tr.Id%2 == 0 {
			list = append(list, tr)
		}
	}

	for _, tr := range list {

		sqlStatement = `
			UPDATE transactions SET status_key = 'CANCELLED' WHERE id = $1;
		`
		cancelTrx, err := tx.Prepare(sqlStatement)
		if err != nil {
			return tx.Rollback()
		}

		_, err = cancelTrx.Exec(tr.Id)
		if err != nil {
			return tx.Rollback()
		}

		sqlStatement = `
			UPDATE account SET account_balance = account_balance - $1 WHERE id = $2;
		`

		updAccountBalance, err := tx.Prepare(sqlStatement)
		if err != nil {
			return tx.Rollback()
		}

		_, err = updAccountBalance.Exec(tr.Amount, accountId)
		if err != nil {
			return tx.Rollback()
		}

		sqlStatement = `
			UPDATE jobs SET last_activity_at = $1 WHERE id = 0;
		`

		updLastActivity, err := tx.Prepare(sqlStatement)
		if err != nil {
			return tx.Rollback()
		}

		s := time.Now().Format(time.RFC3339)
		_, err = updLastActivity.Exec(s)
		if err != nil {
			return tx.Rollback()
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}
