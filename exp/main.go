package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "adam"
	dbname = "picapp"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	// Verify driver name and datasource name are working correctly
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type User struct {
		ID int
		name string
		email string
	}

	var users []User
	rows, err := db.Query(`
		SELECT id, name, email
		FROM users`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.name, &user.email)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		fmt.Println("Row err: ", err)
	}
	fmt.Println(users)
}
