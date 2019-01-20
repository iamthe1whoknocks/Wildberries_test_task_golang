//Logic of working with database
package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//checking error
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//NewDB checks connection to DB from config file and returns sql.DB
func NewDB(datasourse string) (*sql.DB, error) {
	db, err := sql.Open("postgres", datasourse)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}

//struct to parse config file
type dbConfig struct {
	Port     int    `json:"port"`     //порт, на котором слушать запросы
	Endpoint string `json:"endpoint"` //название API
	Host     string `json:"host"`     //hostname, где установлен Postgres
	User     string `json:"user"`     //имя пользователя Postgres
	Password string `json:"password"` //пароль пользователя Postgres
	Schema   string `json:"schema"`   //схема в Postgres
	DBName   string `json:"dbname"`   //название БД в Postgres

}

//GetConnectionString parses config file into struct
func GetConnectionString(filename string) string {
	confFile, err := os.Open("configuration.json")
	checkError(err)

	b, err := ioutil.ReadAll(confFile)
	checkError(err)

	config := &dbConfig{}
	err = json.Unmarshal(b, &config)
	checkError(err)
	connectionStr := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable ", config.User, config.Password, config.Host, config.Port, config.DBName)
	return connectionStr
}

//GetUserComment refers to Postgresql function test.user_comment_get.
func GetUserComment(db *sql.DB, userID int, commentID int) (jsonString []byte, err error) {
	rows, err := db.Query("Select * from test.user_comment_get($1,$2)", userID, commentID)
	if err != nil {
		fmt.Println(err, "GetUserComment error")
		return nil, err
	}
	defer rows.Close()
	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}

//GetUser refers to Postgresql function test.user_get.
func GetUser(db *sql.DB, userID int) (jsonString []byte, err error) {
	rows, err := db.Query("Select * from test.user_get($1)", userID)
	if err != nil {
		fmt.Println(err, "GetUser error")
		return nil, err
	}
	defer rows.Close()
	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}

//GetComment refers to Postgresql function test.comment_get.
func GetComment(db *sql.DB, commentID int) (jsonString []byte, err error) {
	rows, err := db.Query("Select * from test.comment_get($1)", commentID)
	if err != nil {
		fmt.Println(err, "GetComment error")
		return nil, err
	}
	defer rows.Close()
	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}

//PostComment refers to Postgresql function test.user_comment_ins.
func PostComment(db *sql.DB, userID int, body []byte) (jsonString []byte, err error) {

	rows, err := db.Query("Select * from test.user_comment_ins($1,$2)", userID, string(body))
	if err != nil {
		fmt.Println(err, "PostComment error", string(body))
		return nil, err
	}
	defer rows.Close()

	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}

//PutComment refers to Postgresql function test.comment_upd.
func PutComment(db *sql.DB, commentID int, body []byte) (jsonString []byte, err error) {

	rows, err := db.Query("Select * from test.comment_upd($1,$2)", commentID, string(body[:]))
	if err != nil {
		fmt.Println(err, "PutComment error", body)
		return nil, err
	}
	defer rows.Close()

	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}

//DelComment refers to Postgresql function test.comment_del.
func DelComment(db *sql.DB, commentID int) (jsonString []byte, err error) {

	rows, err := db.Query("Select * from test.comment_del($1)", commentID)
	if err != nil {
		fmt.Println(err, "DelComment error")
		return nil, err
	}
	defer rows.Close()

	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}

//DelUser refers to Postgresql function test.user_del.
func DelUser(db *sql.DB, userID int) (jsonString []byte, err error) {

	rows, err := db.Query("Select * from test.user_del($1)", userID)
	if err != nil {
		fmt.Println(err, "DelUser error")
		return nil, err
	}
	defer rows.Close()

	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}

//PostUser refers to Postgresql function test.user_ins.
func PostUser(db *sql.DB, body []byte) (jsonString []byte, err error) {

	rows, err := db.Query("Select * from test.user_ins($1)", body)
	if err != nil {
		fmt.Println(err, "PostUser error", body)
		return nil, err
	}
	defer rows.Close()

	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}

//PutUser refers to Postgresql function test.user_upd.
func PutUser(db *sql.DB, userID int, body []byte) (jsonString []byte, err error) {

	rows, err := db.Query("Select * from test.user_upd($1,$2)", userID, string(body[:]))
	if err != nil {
		fmt.Println(err, "PutUser error", body)
		return nil, err
	}
	defer rows.Close()

	byteStr := make([]byte, 0)
	for rows.Next() {
		err = rows.Scan(&byteStr)
		if err != nil {
			fmt.Println(err, "scanning error")
		}

	}
	return byteStr, nil
}
