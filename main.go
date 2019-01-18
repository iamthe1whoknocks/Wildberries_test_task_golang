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

	if err != nil {
		fmt.Println(err, "exiting programm")
		os.Exit(1)
	}
	myRouter := mux.NewRouter().StrictSlash(true)

	myServer := &Server{
		db:     myDb,
		router: myRouter,
	}
	myServer.Routes()

	server := &http.Server{
		Addr: ":8080",

		Handler: myRouter,
	}
	err = server.ListenAndServe()
	checkError(err)
}

// IndexHandler does...

//curl -d {"txt":"cheers"} -H "Content-Type: application/json" -X POST http://127.0.0.1:8080/api/v1/user/5/comment/
//у некоторых роутов могут быть 2 метода!!!
