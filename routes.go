package main

import (
	"database/sql"

	"github.com/gorilla/mux"
)

//Server does...
type Server struct {
	db     *sql.DB
	router *mux.Router
}

func someFunc() {

}

/*func (s *Server) routes() {
	s.router.HandleFunc("/")
}*/
