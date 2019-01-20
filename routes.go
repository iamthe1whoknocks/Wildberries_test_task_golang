package main

import (
	"database/sql"
	"fmt"
	"golang+postgres/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gorilla/mux"
)

//Server struct
type Server struct {
	db     *sql.DB
	router *mux.Router
}

//Routes is function to avoid big routes in main.go
func (s *Server) Routes() {

	s.router.StrictSlash(true)
	s.router.NotFoundHandler = http.HandlerFunc(notFound)
	s.router.HandleFunc("/", s.IndexHandler())
	s.router.Handle("/metrics", prometheus.Handler())

	subRoute := s.router.PathPrefix("/api/v1/").Subrouter()
	subRoute.HandleFunc("/comment/{commentID:[0-9]+}", s.CommentHandler()).Methods("GET", "PUT", "DELETE")

	subRoute.HandleFunc("/user/{userID:[0-9]+}", s.UserHandler()).Methods("GET", "PUT", "DELETE")

	subRoute.HandleFunc("/user/", s.PostUserHandler()).Methods("POST")

	subRoute.HandleFunc("/user/{userID:[0-9]+}/comment/", s.UserCommentHandler()).Methods("GET", "POST")

	subRoute.HandleFunc("/user/{userID:[0-9]+}/comment/{commentID:[0-9]+}", s.GetUserCommentHandler()).Methods("GET")

	subRoute.HandleFunc("/user/comment/", s.GetAllCommentHandler()).Methods("GET")

}

//IndexHandler is a handler of index page
func (s *Server) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
		duration := time.Since(start)
		reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))
		w.Write([]byte("HELLO!"))

	}
}

//GetUserCommentHandler handles /user/{userID:[0-9]+}/comment/{commentID:[0-9]+} path.
func (s *Server) GetUserCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(405), 405)
			totalReq.WithLabelValues("405", r.URL.String(), r.Method).Inc()
			return
		}
		start := time.Now()
		totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
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
		b, err := models.GetUserComment(s.db, userIDInt, commmentIDInt)
		if err != nil {
			fmt.Println(err, "nothing get")
		}
		if b == nil {
			w.Write([]byte("Cant find comment with such params"))
		}

		duration := time.Since(start)
		reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	}
}

//UserHandler handles /user/{userID:[0-9]+} path
func (s *Server) UserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case "GET":
			start := time.Now()
			totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
			vars := mux.Vars(r)
			userID := vars["userID"]
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				fmt.Println(err)
			}
			b, err := models.GetUser(s.db, userIDInt)
			if err != nil {
				fmt.Println(err, "nothing get")
			}
			if b == nil {
				w.Write([]byte("Cant find user with such params"))
			}
			duration := time.Since(start)
			reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)

		case "DELETE":
			start := time.Now()
			totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
			vars := mux.Vars(r)
			userID := vars["userID"]
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				fmt.Println(err)
			}
			b, err := models.DelUser(s.db, userIDInt)
			duration := time.Since(start)
			reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)

		case "PUT":
			start := time.Now()
			totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
			vars := mux.Vars(r)

			userID := vars["userID"]

			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				fmt.Println(err)
			}
			buf := make([]byte, 0)
			buf, err = ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			b, err := models.PutUser(s.db, userIDInt, buf)

			duration := time.Since(start)
			reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))

			w.Header().Set("Content-Type", "application/json")
			w.Write(b)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(405), 405)
			totalReq.WithLabelValues("405", r.URL.String(), r.Method).Inc()
			return
		}
	}
}

//UserCommentHandler handles /user/{userID:[0-9]+}/comment/ path
func (s *Server) UserCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case "GET":
			start := time.Now()
			totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
			vars := mux.Vars(r)
			userID := vars["userID"]
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				fmt.Println(err)
			}
			commmentIDInt := 0
			b, err := models.GetUserComment(s.db, userIDInt, commmentIDInt)
			if err != nil {
				fmt.Println(err, "nothing get")
			}
			if b == nil {
				w.Write([]byte("Cant find comment with such params"))
			}
			duration := time.Since(start)
			reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		case "POST":
			start := time.Now()
			totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
			vars := mux.Vars(r)
			userID := vars["userID"]
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				fmt.Println(err)
			}
			buf := make([]byte, 0)
			buf, err = ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			b, err := models.PostComment(s.db, userIDInt, buf)
			duration := time.Since(start)
			reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(405), 405)
			totalReq.WithLabelValues("405", r.URL.String(), r.Method).Inc()
			return

		}
	}
}

//GetAllCommentHandler handles /user/comment/ path
func (s *Server) GetAllCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			totalReq.WithLabelValues("405", r.URL.String(), r.Method).Inc()
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(405), 405)
			return
		}
		start := time.Now()
		totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
		b, err := models.GetUserComment(s.db, 0, 0)
		if err != nil {
			fmt.Println(err, "nothing get")
		}
		if b == nil {
			w.Write([]byte("Cant find user with such params"))
		}
		duration := time.Since(start)
		reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

//CommentHandler handles /comment/{commentID:[0-9]+} path.
func (s *Server) CommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case "DELETE":
			start := time.Now()
			totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
			vars := mux.Vars(r)
			commentID := vars["commentID"]
			commentIDInt, err := strconv.Atoi(commentID)
			if err != nil {
				fmt.Println(err)
			}

			b, err := models.DelComment(s.db, commentIDInt)

			duration := time.Since(start)
			reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)

		case "PUT":
			start := time.Now()
			totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
			vars := mux.Vars(r)
			commentID := vars["commentID"]
			commentIDInt, err := strconv.Atoi(commentID)
			if err != nil {
				fmt.Println(err)
			}
			buf := make([]byte, 0)
			buf, err = ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			b, err := models.PutComment(s.db, commentIDInt, buf)

			duration := time.Since(start)
			reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))

			w.Header().Set("Content-Type", "application/json")
			w.Write(b)

		case "GET":
			start := time.Now()
			totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
			vars := mux.Vars(r)
			commentID := vars["commentID"]
			commentIDInt, err := strconv.Atoi(commentID)
			if err != nil {
				fmt.Println(err)
			}
			b, err := models.GetComment(s.db, commentIDInt)
			if err != nil {
				fmt.Println(err, "nothing get")
			}
			if b == nil {
				w.Write([]byte("Cant find user with such params"))
			}

			duration := time.Since(start)
			reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))

			w.Header().Set("Content-Type", "application/json")
			w.Write(b)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			http.Error(w, http.StatusText(405), 405)
			totalReq.WithLabelValues("405", r.URL.String(), r.Method).Inc()
			return
		}
	}
}

//PostUserHandler handles /user/ path.
func (s *Server) PostUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			totalReq.WithLabelValues("405", r.URL.String(), r.Method).Inc()
			return
		}
		start := time.Now()
		totalReq.WithLabelValues("200", r.URL.String(), r.Method).Inc()
		buf := make([]byte, 0)
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		defer r.Body.Close()

		b, err := models.PostUser(s.db, buf)

		duration := time.Since(start)
		reqDuration.WithLabelValues("200", r.URL.String(), r.Method).Observe(float64(duration.Nanoseconds()))

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	}
}

//custom 404 page,it was made to collect metrics from 404 page
func notFound(w http.ResponseWriter, r *http.Request) {
	totalReq.WithLabelValues("404", r.URL.String(), r.Method).Inc()
	fmt.Fprint(w, "Custom 404 Page not found")
}
