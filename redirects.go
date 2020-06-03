package main

import (
	"net/http"
)

//Home is the main function
func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/home.html")
		return
	}
}

//ListUsers list all users in db
func ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/listUsers.html")
		return
	}
}
