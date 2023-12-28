package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/gogrpc"))
	if err != nil {
		fmt.Println("Error connecting to database")
		panic(err)
	}

	fmt.Println("Database is now connected")
 return db
}