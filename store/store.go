package store

import (
	"github.com/ashtanko/octo-server/model"
	"time"
)

type Store interface {
	Account() AccountStore
	Transaction() TransactionStore
	Job() JobStore
	Close()
}

type AccountStore interface {
	Create(*model.Account) error
	Find(id string) (*model.Account, error)
	GetBalance(id string) (float64, error)
	SetBalance(id string, balance float64) error
	UpdateBalance(id string, balance float64) error
	Delete(id string) error
}

type TransactionStore interface {
	Save(*model.IncomingRequest, string) error
	Fetch(int) ([]model.Transaction, error)
	Update(int) error
}

type JobStore interface {
	CancelTransactionAndCorrectBalance(jobRunInterval time.Duration, accountId int) error
}
