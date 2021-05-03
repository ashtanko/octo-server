package app

import (
	"encoding/json"
	"fmt"
	"github.com/ashtanko/octo-server/model"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"net/http"
)

func (s *Server) InitAccount() {
	s.RootRouter.HandleFunc("/account/create", s.createAccount()).Methods("POST")
	s.RootRouter.HandleFunc("/account/{account_id:[0-9]+}", s.getAccount()).Methods("GET")
	s.RootRouter.HandleFunc("/account/balance", s.getAccountBalance()).Methods("GET")
	s.RootRouter.HandleFunc("/account/{account_id:[0-9]+}", s.deleteAccount()).Methods("DELETE")
}

func (s *Server) getAccountBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		u, err := s.Store.Account().GetBalance("0")
		if err != nil {
			s.error(w, http.StatusNotFound, err)
			return
		}
		s.respond(w, http.StatusOK, u)
	}
}

func (s *Server) deleteAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["account_id"]
		err := s.Store.Account().Delete(id)
		if err != nil {
			s.error(w, http.StatusNotFound, err)
			return
		}
		s.respond(w, http.StatusOK, fmt.Sprintf("Account with id %v deleted", id))
	}
}

func (s *Server) getAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["account_id"]
		u, err := s.Store.Account().Find(id)
		if err != nil {
			s.error(w, http.StatusNotFound, err)
			return
		}
		s.respond(w, http.StatusOK, u)
	}
}

func (s *Server) createAccount() http.HandlerFunc {

	type request struct {
		AccountBalance float64 `json:"account_balance"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		u := &model.Account{
			AccountBalance: req.AccountBalance,
		}

		if err := s.Store.Account().Create(u); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, http.StatusCreated, u)
	}
}

func (s *Server) updateBalance(id string, balance float64) {
	if err := s.Store.Account().UpdateBalance(id, balance); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == (checkViolationErrorCode) {
				err := s.Store.Account().SetBalance(id, 0)
				if err != nil {
					s.Logger.Error(err)
				}
				return
			}
		}
	}
	s.Logger.Printf("Balance updated %v %v", id, balance)
}
