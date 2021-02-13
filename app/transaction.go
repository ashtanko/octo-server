package app

import (
	"encoding/json"
	"github.com/ashtanko/octo-server/model"
	"github.com/lib/pq"
	"net/http"
)

const (
	stateWin                 = "win"
	stateLost                = "lost"
	uniqueViolationErrorCode = "23505"
	transactionStatusDone    = "DONE"
)

func (s *Server) InitTransaction() {
	s.RootRouter.HandleFunc("/transaction", s.handTransaction()).Methods("POST")
	s.RootRouter.HandleFunc("/transaction", s.getTransactions()).Methods("GET")
}

func (s *Server) getTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := s.Store.Transaction().Fetch(20)
		if err != nil {
			s.respond(w, http.StatusOK, list)
		}
		s.respond(w, http.StatusOK, list)
	}
}

func (s *Server) handTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		source := r.Header.Get("Source-Type")

		req := &model.IncomingRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		income := &model.IncomingRequest{
			Source:        source,
			State:         req.State,
			Amount:        req.Amount,
			TransactionID: req.TransactionID,
		}

		switch income.State {
		case stateLost:
			income.Amount = -income.Amount
		case stateWin:
		// ignore
		default:
			s.respond(w, http.StatusCreated, income)
			return
		}

		if err := s.Store.Transaction().Save(income, transactionStatusDone); err != nil {
			if err, ok := err.(*pq.Error); ok {
				s.Logger.Error(err.Code.Name())
				if err.Code == (uniqueViolationErrorCode) {
					s.respond(w, http.StatusCreated, income)
					return
				}
			}
		}

		s.updateBalance(model.StaticUserId, income.Amount)

		s.respond(w, http.StatusCreated, income)
	}
}
