package main

//User is the model to serializee
type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
}
