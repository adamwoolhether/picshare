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

/*	// Method 1
	newDB := db.Where("email = ?", "blah@blah.com").First(&u)
	if newDB.Error != nil {
		panic(newDB.Error)
	}*/

/*	// Method 2
	db = db.Where("name = ?", "frank").First(&u)
	errors := db.Error
	if db.Error != nil {
		fmt.Println(errors)
	}*/

	// Method 3
	if err := db.Where("name = ?", "frank").First(&u).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			fmt.Println("Record is not found")
		default:
			panic(err)
		}
	}


	//db.First(&u)
	fmt.Println(u)

}
