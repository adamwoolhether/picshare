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
	Orders []Order
}

type Order struct {
	gorm.Model
	UserID	uint
	Amount 	int
	Description string
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
	db.AutoMigrate(&User{}, &Order{})

	var u User
	// Make sure to preload the orders
	if err := db.Preload("Orders").First(&u).Error; err != nil {
		panic(err)
	}
	//createOrder(db, u, 1001, "Fake Description #1")
	//createOrder(db, u, 9999, "Fake Description #2")
	//createOrder(db, u, 100, "Fake Description #3")
	fmt.Println(u)
	fmt.Println(u.Orders)

	// Another method:
	var users []User
	if err = db.Preload("Orders").Find(&users).Error; err != nil {
		panic(err)
	}
	fmt.Println(users)

}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	err := db.Create(&Order{
		UserID: user.ID,
		Amount: amount,
		Description: desc,
	}).Error
	if err != nil {
		panic(err)
	}
}