package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/MeguMan/userBalanceAvitoTask/internal/model"
	"github.com/MeguMan/userBalanceAvitoTask/internal/store/postgres_store"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	s.router.HandleFunc("/createuser", s.CreateUser()).Methods("POST")
}

func (s *server) CreateUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := model.User{}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ur := s.store.User()

		_, err = ur.Create(&u)
		if err != nil {
			log.Print(err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, u)
	}
}

