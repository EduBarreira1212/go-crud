package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Conect()
	if err != nil {
		w.Write([]byte("Error to connect with database!"))
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from user")
	if err != nil {
		w.Write([]byte("Error to select from database!"))
		return
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var user user

		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Error to scan from database!"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Error to convert users to JSON"))
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error to parse user id!"))
		return
	}

	db, err := database.Conect()
	if err != nil {
		w.Write([]byte("Error to connect with database!"))
		return
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM user WHERE id = ?", ID)
	if err != nil {
		w.Write([]byte("Error to get user!"))
		return
	}
	defer row.Close()

	var user user
	if row.Next() {
		err := row.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			w.Write([]byte("Error to scan user!"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Error to convert user to JSON!"))
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	ID, err := strconv.ParseUint(parameters["id"], 10, 32)
	if err != nil {
		w.Write([]byte("Error to parse user id!"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error to read body!"))
		return
	}

	var user user
	if err := json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Error to convert user to struct"))
		return
	}

	db, err := database.Conect()
	if err != nil {
		w.Write([]byte("Error to connect with database!"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE user SET name = ?, email = ? where id = ?")
	if err != nil {
		w.Write([]byte("Error to create statement"))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(user.Name, user.Email, ID); err != nil {
		w.Write([]byte("Error to execute statement"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
