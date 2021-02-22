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
	us.DestructiveReset()
	us.AutoMigrate()
	user := models.User{
		Name:     "adam",
		Email:    "adam@wade.com",
		Password: "adam",
		Remember: "abc123",
	}
	err = us.Create(&user); if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", user)

	user2, err := us.ByRemember("abc123")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *user2)
}
