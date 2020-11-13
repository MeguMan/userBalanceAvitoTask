package apiserver

import (
	"errors"
	"fmt"
	"github.com/MeguMan/userBalanceAvitoTask/internal/model"
	"github.com/MeguMan/userBalanceAvitoTask/internal/store/postgres_store"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"regexp"
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
	s.router.HandleFunc("/users/{id}/balance", s.GetUserBalance()).Methods("GET")
	s.router.HandleFunc("/users/{id}/balance/add", s.IncreaseUserBalance()).Methods("PUT")
	s.router.HandleFunc("/users/{id}/balance/reduce", s.ReduceUserBalance()).Methods("PUT")
	s.router.HandleFunc("/users/{sender_id}/balance/transfer", s.TransferUserBalance()).Methods("PUT")
}

func (s *server) GetUserBalance() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		u := model.User{}
		vars := mux.Vars(r)
		u.ID, _ = strconv.Atoi(vars["id"])

		currencyName := r.URL.Query().Get("currency")

		ur := s.store.User()
		var err error
		u.Balance, err = ur.GetBalanceById(u.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if currencyName != "" {
			u.Balance = ConvertCurrency(u.Balance, currencyName)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, u.Balance)
	}
}

func (s *server) IncreaseUserBalance() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		u := model.User{}
		vars := mux.Vars(r)
		balanceToAdd := r.URL.Query().Get("balance")
		if balanceToAdd == "" {
			err :=  errors.New("parameter balance is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.ID, _ = strconv.Atoi(vars["id"])
		u.Balance, _ = strconv.ParseFloat(balanceToAdd, 64)

		ur := s.store.User()
		err := ur.IncreaseBalance(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Balance was increased")
	}
}

func (s *server) ReduceUserBalance() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		u := model.User{}
		vars := mux.Vars(r)
		balanceToWriteOff:= r.URL.Query().Get("balance")
		if balanceToWriteOff == "" {
			err :=  errors.New("parameter balance is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.ID, _ = strconv.Atoi(vars["id"])
		u.Balance, _ = strconv.ParseFloat(balanceToWriteOff, 64)

		ur := s.store.User()
		err := ur.ReduceBalance(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Balance was reduced")
	}
}

func (s *server) TransferUserBalance() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sender := model.User{}
		receiver := model.User{}

		receiverId := r.URL.Query().Get("receiver_id")
		if receiverId == "" {
			err :=  errors.New("parameter receiver_id is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		balanceForTransfer := r.URL.Query().Get("balance")
		if balanceForTransfer == "" {
			err :=  errors.New("parameter balance is empty")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		sender.ID, _ = strconv.Atoi(vars["sender_id"])
		sender.Balance, _ = strconv.ParseFloat(balanceForTransfer, 64)
		receiver.ID, _ = strconv.Atoi(receiverId)
		receiver.Balance, _ = strconv.ParseFloat(balanceForTransfer, 64)

		ur := s.store.User()
		err := ur.TransferBalance(sender, receiver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Balance was transfered")
	}
}

func ConvertCurrency(balance float64, currencyName string) float64 {
	req, _ := http.Get(fmt.Sprintf("https://api.exchangeratesapi.io/latest?symbols=RUB&base=%s", currencyName))
	body, _ := ioutil.ReadAll(req.Body)
	re := regexp.MustCompile(`(\d\d.\d*)`)
	data := re.FindAllString(string(body), -1)
	exchangeRate, _ := strconv.ParseFloat(data[0], 64)
	balance = balance / exchangeRate
	return balance
}
