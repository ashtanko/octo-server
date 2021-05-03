package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ashtanko/octo-server/app/config"
	"github.com/ashtanko/octo-server/model"
	"github.com/ashtanko/octo-server/store"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	DBPingAttempts    = 18
	DBPingTimeoutSecs = 10
	ExitDBOpen        = 101
	ExitPing          = 102
)

// Store ...
type Store struct {
	Db               *sql.DB
	accountStore     *AccountStore
	dataSource       string
	transactionStore *TransactionStore
	jobStore         *JobStore
}

// CreateNewAndConnect ...
func CreateNewAndConnect(cfg config.Config) *Store {
	dataSourceName := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Database,
		cfg.DB.Port)

	s := &Store{
		dataSource: dataSourceName,
	}
	s.initDBConnection()
	return s
}

func (s *Store) initDBConnection() {
	s.Db = setupDBConnection(s.dataSource)
}

func setupDBConnection(dataSourceName string) *sql.DB {
	db, err := sql.Open(model.DatabaseDriverPostgres, dataSourceName)
	if err != nil {
		logrus.Fatal("Failed to open SQL connection ", err)
		time.Sleep(time.Second)
		os.Exit(ExitDBOpen)
	}
	for i := 0; i < DBPingAttempts; i++ {
		logrus.Info("Pinging SQL: ", dataSourceName)
		ctx, cancel := context.WithTimeout(context.Background(), DBPingTimeoutSecs*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err == nil {
			break
		} else {
			if i == DBPingAttempts-1 {
				logrus.Fatal("Failed to ping DB, server will exit ", err)
				time.Sleep(time.Second)
				os.Exit(ExitPing)
			} else {
				logrus.Error("Failed to ping DB ", err)
				time.Sleep(DBPingTimeoutSecs * time.Second)
			}
		}
	}
	return db
}

// Close db connection
func (s *Store) Close() {
	err := s.Db.Close()
	if err != nil {
		logrus.Error("Failed to close db")
	}
}

// Transaction sql store
func (s *Store) Transaction() store.TransactionStore {
	if s.transactionStore != nil {
		return s.transactionStore
	}
	s.transactionStore = &TransactionStore{
		store: s,
	}

	return s.transactionStore
}

// Account sql store
func (s *Store) Account() store.AccountStore {
	if s.accountStore != nil {
		return s.accountStore
	}

	s.accountStore = &AccountStore{
		Store: s,
	}

	return s.accountStore
}

// Job sql store
func (s *Store) Job() store.JobStore {
	if s.jobStore != nil {
		return s.jobStore
	}

	s.jobStore = &JobStore{
		Store: s,
	}

	return s.jobStore
}
