package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "adam"
	dbname   = "picapp"
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

	//// Check connection to db
	//err = db.Ping()
	//if err != nil {
	//	panic(err)
	//}

/*	var id int
	err = db.QueryRow(`
	INSERT INTO users(name, email)
	VALUES($1, $2)
	RETURNING id`, "Joe Wade", "joe@gmail.com").Scan(&id)
	if err != nil {
		panic(err)
	}*/

	// If you don't want to chain the command, run the above code block as follows
	var id int
	var name, email string
	row := db.QueryRow(`
		SELECT id, name, email
		FROM users
		WHERE id=$1`, 6)
	err = row.Scan(&id, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows with given id")
		} else {
			panic(err)
		}
	}

	fmt.Println("id: ", id, "name: ", name, "email: ", email)
}