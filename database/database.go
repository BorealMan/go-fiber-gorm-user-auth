package database

import (
	"app/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB() {

	var err error
	// DB, err = gorm.Connect("mysql", GenerateDBURL())
	DB, err = gorm.Open(mysql.Open(GenerateDBURL()), &gorm.Config{})

	if err != nil {
		log.Fatalf(`\nFailed To Connect To Database: %v+\n`, err)
	}

	fmt.Println("\nSuccessfully Connected To Database")
	// seed.Seed()
}

func GenerateDBURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		config.DB_USERNAME,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_DATABASE,
	)
}
