package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/phildehovre/go-complete-api/services/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}
func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("api/v1").Subrouter()

	store := user.NewStore(s.db)

	userHandler := user.NewHandler(store)
	userHandler.RegisterRoutes(subrouter)
	

	log.Println("listening on ", s.addr)
	return http.ListenAndServe(s.addr, router)
}
