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

func (s *Server) Routes (){
	s.router.HandleFunc("/", IndexHandler)

	subRoute := s.router.PathPrefix("/api/v1/").Subrouter()
	subRoute.HandleFunc("/comment/{commentID:[0-9]+}", s.CommentHandler()).Methods("GET", "PUT", "DELETE")

	subRoute.HandleFunc("/user/{userID:[0-9]+}", s.UserHandler()).Methods("GET", "PUT", "DELETE")

	subRoute.HandleFunc("/user/", s.PostUserHandler()).Methods("POST")

	subRoute.HandleFunc("/user/{userID:[0-9]+}/comment/", s.UserCommentHandler()).Methods("GET", "POST")

	subRoute.HandleFunc("/user/{userID:[0-9]+}/comment/{commentID:[0-9]+}", s.GetUserCommentHandler()).Methods("GET")

	subRoute.HandleFunc("/user/comment/", s.GetAllCommentHandler()).Methods("GET")

}


func (s *Server) IndexHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!"))
}
}

//GetUserCommentHandler does...
func (s *Server) GetUserCommentHandler() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
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
}

//UserHandler does...
func (s *Server) UserHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
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
}

//UserCommentHandler does...
func (s *Server) UserCommentHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
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
}

//GetAllCommentHandler does...
func (s *Server) GetAllCommentHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
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
}

//CommentHandler does....
func (s *Server) CommentHandler() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
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
}

//DelUserHandler does....

//PostUserHandler does...
func (s *Server) PostUserHandler() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
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
}