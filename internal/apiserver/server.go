package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/MeguMan/userBalanceAvitoTask/internal/model"
	"github.com/MeguMan/userBalanceAvitoTask/internal/store/postgres_store"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	router *mux.Router
	store  postgres_store.Store
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer(store postgres_store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users/{id}/balance/add", s.AddUserBalance()).Methods("PUT")
	s.router.HandleFunc("/users/{id}/balance/reduce", s.ReduceUserBalance()).Methods("PUT")
	s.router.HandleFunc("/users/{id}/balance", s.GetUserBalance()).Methods("GET")
}

func (s *server) GetUserBalance() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := model.User{}
		vars := mux.Vars(r)
		u.ID, _ = strconv.Atoi(vars["id"])

		w.Header().Set("Content-Type", "application/json")

		ur := s.store.User()

		var err error
		u.Balance, err = ur.GetBalanceById(u)
		if err != nil {
			log.Print(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, u)
	}
}

func (s *server) AddUserBalance() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := model.User{}
		vars := mux.Vars(r)
		u.ID, _ = strconv.Atoi(vars["id"])

		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ur := s.store.User()
		err = ur.AddBalance(u)
		if err != nil {
			log.Print(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Balance was increased")
	}
}

func (s *server) ReduceUserBalance() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := model.User{}
		vars := mux.Vars(r)
		u.ID, _ = strconv.Atoi(vars["id"])

		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ur := s.store.User()
		err = ur.ReduceBalance(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Balance was decreased")
	}
}

