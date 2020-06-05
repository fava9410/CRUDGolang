package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

func main() {

	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	http.HandleFunc("/", Home)
	http.HandleFunc("/listUsers", ListUsers)

	http.HandleFunc("/createUser", CreateUser)
	http.HandleFunc("/getUsers", GetUsers)
	http.HandleFunc("/getUserById", GetUser)
	http.HandleFunc("/updateUser", UpdateUser)
	http.HandleFunc("/deleteUser", DeleteUser)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
