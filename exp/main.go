package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "adam"
	dbname = "picapp"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
	Color string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	// Verify driver name and datasource name are working correctly
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	//db.Migrator().DropTable(&User{})
	db.AutoMigrate(&User{})

	var u User
	/*//newDB := db.Where("id = ? AND color = ?", 7, "green") //GORM uses a '?' as an arg placeholder instead of $1
	//newDB.First(&u)*/

/*	// Another way to chain queries:
 	db.Where("color = ?", "green").
		Where("id > ?", 3).
		First(&u)*/

/*	// Or you can alter the user object:
	var u User = User{
		Color: "green",
		Email: "jane@jane.com",
	}
	db.Where(&u).First(&u)*/

/*	// Using Find example:
	var users []User
	db.Find(&users)
	fmt.Println(len(users))
	fmt.Println(users)*/


	db.First(&u)
	fmt.Println(u)

}
