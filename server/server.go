package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error to read body!"))
		return
	}

	var user user
	if err = json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Error to convert user to struct"))
		return
	}

	db, err := database.Conect()
	if err != nil {
		w.Write([]byte("Error to conect to database"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into user (name, email) values (?, ?)")
	if err != nil {
		w.Write([]byte("Error to create statement"))
		return
	}
	defer statement.Close()

	insertion, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.Write([]byte("Error to execute statement"))
		return
	}

	idInserted, err := insertion.LastInsertId()
	if err != nil {
		w.Write([]byte("Error to get id inserted"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User created with sucess! Id: %d", idInserted)))
}
