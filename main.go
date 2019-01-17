// golang+postgres project main.go
package main

import (
	"database/sql"
	"fmt"

	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"golang+postgres/models"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

//Env does...
type Env struct {
	Db *sql.DB
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	db, err := models.NewDB(models.GetConnectionString("config.json"))
	if err != nil {
		fmt.Println(err, "exiting programm")
		os.Exit(1)
	}
	fmt.Println("Connected to DB")

	env := &Env{Db: db}

	if err != nil {
		fmt.Println(err, "exiting programm")
		os.Exit(1)
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", IndexHandler)

	subRoute := router.PathPrefix("/api/v1/").Subrouter()
	subRoute.HandleFunc("/comment/{commentID:[0-9]+}", env.CommentHandler).Methods("GET", "PUT", "DELETE")

	subRoute.HandleFunc("/user/{userID:[0-9]+}", env.UserHandler).Methods("GET", "PUT", "DELETE")

	subRoute.HandleFunc("/user/", env.PostUserHandler).Methods("POST")

	subRoute.HandleFunc("/user/{userID:[0-9]+}/comment/", env.UserCommentHandler).Methods("GET", "POST")

	subRoute.HandleFunc("/user/{userID:[0-9]+}/comment/{commentID:[0-9]+}", env.GetUserCommentHandler).Methods("GET")

	subRoute.HandleFunc("/user/comment/", env.GetAllCommentHandler).Methods("GET")

	server := &http.Server{
		Addr: ":8080",

		Handler: router,
	}
	err = server.ListenAndServe()
	checkError(err)
}

// IndexHandler does...
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!"))
}

//GetUserCommentHandler does...
func (env *Env) GetUserCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	vars := mux.Vars(r)
	userID := vars["userID"]
	commentID := vars["commentID"]

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println(err)
	}

	commmentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		fmt.Println(err)
	}
	b, err := models.GetUserComment(env.Db, userIDInt, commmentIDInt)
	if err != nil {
		fmt.Println(err, "nothing get")
	}
	if b == nil {
		w.Write([]byte("Cant find comment with such params"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

//UserHandler does...
func (env *Env) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		vars := mux.Vars(r)
		userID := vars["userID"]
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println(err)
		}
		b, err := models.GetUser(env.Db, userIDInt)
		if err != nil {
			fmt.Println(err, "nothing get")
		}
		if b == nil {
			w.Write([]byte("Cant find user with such params"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	case "DELETE":
		vars := mux.Vars(r)
		userID := vars["userID"]
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println(err)
		}
		b, err := models.DelUser(env.Db, userIDInt)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	case "PUT":
		vars := mux.Vars(r)

		userID := vars["userID"]

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println(err)
		}
		buf := make([]byte, 0)
		buf, err = ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		b, err := models.PutUser(env.Db, userIDInt, buf)

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(405), 405)
		return
	}
}

//UserCommentHandler does...
func (env *Env) UserCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		vars := mux.Vars(r)
		userID := vars["userID"]
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println(err)
		}
		commmentIDInt := 0
		b, err := models.GetUserComment(env.Db, userIDInt, commmentIDInt)
		if err != nil {
			fmt.Println(err, "nothing get")
		}
		if b == nil {
			w.Write([]byte("Cant find comment with such params"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case "POST":
		vars := mux.Vars(r)
		userID := vars["userID"]
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			fmt.Println(err)
		}
		buf := make([]byte, 0)
		buf, err = ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		b, err := models.PostComment(env.Db, userIDInt, buf)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(405), 405)
		return

	}
}

//GetAllCommentHandler does...
func (env *Env) GetAllCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	b, err := models.GetUserComment(env.Db, 0, 0)
	if err != nil {
		fmt.Println(err, "nothing get")
	}
	if b == nil {
		w.Write([]byte("Cant find user with such params"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

//CommentHandler does....
func (env *Env) CommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "DELETE":
		vars := mux.Vars(r)
		commentID := vars["commentID"]
		commentIDInt, err := strconv.Atoi(commentID)
		if err != nil {
			fmt.Println(err)
		}
		b, err := models.DelComment(env.Db, commentIDInt)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	case "PUT":
		vars := mux.Vars(r)
		commentID := vars["commentID"]
		commentIDInt, err := strconv.Atoi(commentID)
		if err != nil {
			fmt.Println(err)
		}
		buf := make([]byte, 0)
		buf, err = ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		b, err := models.PutComment(env.Db, commentIDInt, buf)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	case "GET":
		vars := mux.Vars(r)
		commentID := vars["commentID"]
		commentIDInt, err := strconv.Atoi(commentID)
		if err != nil {
			fmt.Println(err)
		}
		b, err := models.GetComment(env.Db, commentIDInt)
		if err != nil {
			fmt.Println(err, "nothing get")
		}
		if b == nil {
			w.Write([]byte("Cant find user with such params"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(405), 405)
		return
	}
}

//DelUserHandler does....

//PostUserHandler does...
func (env *Env) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	buf := make([]byte, 0)
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()

	b, err := models.PostUser(env.Db, buf)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}

//curl -d {"txt":"cheers"} -H "Content-Type: application/json" -X POST http://127.0.0.1:8080/api/v1/user/5/comment/
//у некоторых роутов могут быть 2 метода!!!
