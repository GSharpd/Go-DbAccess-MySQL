package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Insert a User into the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read request body"))
		return
	}

	var user user

	if err = json.Unmarshal(requestBody, &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to convert request body to user: " + err.Error()))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to connect to database: " + err.Error()))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into users (name, email) values (?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create sql statement: " + err.Error()))
		return
	}
	defer statement.Close()

	insert, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to execute statement: " + err.Error()))
		return
	}

	lastInsertedId, err := insert.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to obtain id: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User inserted with ID: %d", lastInsertedId)))
}

// Get all Users from database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to connect to database: " + err.Error()))
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get all users: " + err.Error()))
		return
	}
	defer rows.Close()

	var users []user

	for rows.Next() {
		var user user

		if err = rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to parse returned object into user struct: " + err.Error()))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode users | " + err.Error()))
		return
	}
}

// Get one user from database
func GetUser(w http.ResponseWriter, r *http.Request) {

}

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
