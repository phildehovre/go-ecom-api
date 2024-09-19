package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/phildehovre/go-complete-api/services/cart"
	"github.com/phildehovre/go-complete-api/services/order"
	"github.com/phildehovre/go-complete-api/services/product"
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
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	store := user.NewStore(s.db)
	userHandler := user.NewHandler(store)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore, productStore)
	cartHandler.RegisterRoutes(subrouter)

	return http.ListenAndServe(s.addr, router)
}
