package api

import (
	"database/sql"
	"net/http"

	"github.com/nagy-gergely/api-test/services/user"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{address: address, db: db}
}

func (server *APIServer) Start() error {
	router := http.NewServeMux()
	apiRouter := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	userStore := user.NewStore(server.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(apiRouter)

	return http.ListenAndServe(server.address, router)
}
