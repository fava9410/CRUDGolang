package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

var testdb *sql.DB

func init() {
	fmt.Println("entre al init test")
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test")

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Fatalf(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	//defer db.Close()
	//para test no se puede cerrar la conexion
}

func TestGetUser(t *testing.T) {

	req, err := http.NewRequest("GET", "/getUserById", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "1")

	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUser)

	handler.ServeHTTP(rr, req)

	//falta test para saber si el json que retorna es correcto
	/*
		rawJSON := `{"id":"1","firstname":"hola","lastname":"mundo","email":"hhh@cl.com","gender":"1"}`
		responseJSON := rr.Body.String()

		fmt.Println(rawJSON)
		fmt.Println(responseJSON)
		fmt.Println(rawJSON == responseJSON)

		if rawJSON != responseJSON {
			t.Errorf("error, got '%v' and receive '%v'", rr.Body.String(), rawJSON)
		}*/

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("error en el content-type")
	}

}

func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/getUsers", nil)

	if err != nil {
		t.Errorf("error in new request")
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(GetUsers)
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("content type")
	}

	JSONUsers := rr.Body.String()
	var users []User
	json.Unmarshal([]byte(JSONUsers), &users)

	var totalusers int
	var totalusersJSON int

	totalusersJSON = len(users)

	_ = db.QueryRow("select count(*) from user;").Scan(&totalusers)

	if totalusersJSON != totalusers {
		t.Errorf("diferente cantidad de usuarios %d y %d", totalusersJSON, totalusers)
	}

}

func TestCreateUser(t *testing.T) {
	var userTest User
	userTest.Firstname = "nombre_testing"
	userTest.Lastname = "apellido_testing"
	userTest.Email = "hola_testing@testinggo.com"
	userTest.Gender = "2"

	data := url.Values{}
	data.Add("first_name", userTest.Firstname)
	data.Add("last_name", userTest.Lastname)
	data.Add("email", userTest.Email)
	data.Add("gender", userTest.Gender)

	req, err := http.NewRequest("POST", "/createUser", strings.NewReader(data.Encode()))

	if err != nil {
		t.Errorf("error in new request")
	}

	//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CreateUser)
	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 301 {
		t.Errorf("Error in status code, got %d", rr.Result().StatusCode)
	}

	row := db.QueryRow("select * from user where email = ? and lastname = ?;", userTest.Email, userTest.Lastname)

	var userTestDB User
	switch err := row.Scan(&userTestDB.ID, &userTestDB.Firstname, &userTestDB.Lastname, &userTestDB.Email, &userTestDB.Gender); err {
	case sql.ErrNoRows:
		t.Errorf("Insert failed")
	case nil:
		userTestDB.ID = ""
		if userTest != userTestDB {
			t.Errorf("They are different, inserted %v , select %v", userTest, userTestDB)
		}
	default:
		t.Errorf("Algo fallo leyendo para comparar %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	data := url.Values{}
	data.Add("first_name", "hola_testing")
	data.Add("last_name", "hola_testing")
	data.Add("email", "hola_testing@testing.com")
	data.Add("gender", "1")
	data.Add("id", "333")

	req, err := http.NewRequest("POST", "/updateUser", strings.NewReader(data.Encode()))

	if err != nil {
		t.Errorf("error in new request")
	}

	//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(UpdateUser)

	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 202 {
		t.Errorf("status code is diferent, got %d", rr.Result().StatusCode)
	}
}

func TestDeleteUser(t *testing.T) {
	data := url.Values{}
	data.Set("id", "12")

	req, err := http.NewRequest("POST", "/deleteUser", strings.NewReader(data.Encode()))

	if err != nil {
		t.Errorf("Error new request in delete user")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(DeleteUser)
	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 301 {
		t.Errorf("error en delete, %v", rr.Result().StatusCode)
	}
}

func TestMethodInDelete(t *testing.T) {
	req, _ := http.NewRequest("GET", "/deleteUser", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(DeleteUser)
	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Delete problem with method %v", rr.Result().StatusCode)
	}
}
