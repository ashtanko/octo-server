package mockstore

import (
	"time"
)

type JobStore struct {
	store *Store
}

func (j *JobStore) CancelTransactionAndCorrectBalance(jobRunInterval time.Duration, accountId int) error {
	return nil
}
