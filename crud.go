package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

//CreateUser gets data from form and create user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user User

	user.Firstname = r.PostFormValue("first_name")
	user.Lastname = r.FormValue("last_name")
	user.Email = r.FormValue("email")
	user.Gender = r.FormValue("gender")

	// Save to database
	smt, err := db.Prepare("Insert into user(firstname, lastname, email, gender) values(?,?,?,?)")

	if err != nil {
		fmt.Println("Prepare error")
		panic(err)
	}

	_, err = smt.Exec(user.Firstname, user.Lastname, user.Email, user.Gender)

	if err != nil {
		fmt.Println("me volvi a romper")
		panic(err)
	}

	http.Redirect(w, r, "/listUsers", 301)

}

//GetUsers returns all users in db
func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entre get users")
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	rows, err := db.Query("select * from user;")

	if err != nil {
		fmt.Println("error en el select")
	}

	var users []User
	var user User
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Gender)
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

//GetUser gets user by id
func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get user by id")
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	userID := r.FormValue("id")

	row := db.QueryRow("select * from user where id = ?;", userID)

	var user User
	switch err := row.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Gender); err {
	case sql.ErrNoRows:
		fmt.Println("User doesnt exist")
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	default:
		panic(err)
	}

}

//UpdateUser updates an specific user with data from frontend
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var user User

	user.Firstname = r.PostFormValue("first_name")
	user.Lastname = r.FormValue("last_name")
	user.Email = r.FormValue("email")
	user.Gender = r.FormValue("gender")
	user.ID = r.FormValue("id")

	stm, err := db.Prepare("Update user set firstname=?, lastname=?, email=?, gender=? where id=?")

	if err != nil {
		fmt.Println("error preparando query")
	}

	_, err = stm.Exec(user.Firstname, user.Lastname, user.Email, user.Gender, user.ID)

	if err != nil {
		fmt.Println("error in update")
		http.Error(w, "error in update", http.StatusInternalServerError)
	}

	//http.Redirect(w, r, "/listUsers", 301)
	w.WriteHeader(http.StatusAccepted)
}

//DeleteUser deletes an user with its id
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	userID := r.FormValue("id")

	res, err := db.Exec("delete from user where id = ?;", userID)

	if err != nil {
		fmt.Println("error eliminando usuario")
	}

	totalRegisterAffected, err := res.RowsAffected()

	if totalRegisterAffected != 1 {
		http.Error(w, "total registers affected were different than 1", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}
