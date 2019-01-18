// golang+postgres project main.go
package main

import (
	"fmt"

	"net/http"
	"os"

	"golang+postgres/models"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	myDb, err := models.NewDB(models.GetConnectionString("config.json"))
	if err != nil {
		fmt.Println(err, "exiting programm")
		os.Exit(1)
	}

	fmt.Println("Connected to DB")

	myRouter := mux.NewRouter()

	myServer := &Server{
		db:     myDb,
		router: myRouter,
	}

	//function with all routes
	myServer.Routes()

	server := &http.Server{
		Addr: ":8080",

		Handler: myRouter,
	}
	err = server.ListenAndServe()
	checkError(err)
}

//curl -d {"txt":"cheers"} -H "Content-Type: application/json" -X POST http://127.0.0.1:8080/api/v1/user/5/comment
