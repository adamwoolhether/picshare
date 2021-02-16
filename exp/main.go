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

	user := models.User{
		Name: "nikki liao",
		Email: "nikki@cat.com",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}
	//user.Email = "jon@newemail.com"
	//if err := us.Update(&user); err != nil {
	//	panic(err)
	//}
	//userByEmail, err := us.ByEmail("jon@newemail.com")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(userByEmail)
	if err := us.Delete(user.ID); err != nil {
		panic(err)
	}
	userByID, err := us.ByID(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(userByID)
}