package main

import (
	"fmt"
	"picapp/models"
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
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	//us.DestructiveReset()
	user, err := us.ByID(2)
	fmt.Println(user)
	if err != nil {
		panic(err)
	}
}