package mysql_driver

import (
	"fmt"
	"log"
	"os"

	"backend/drivers/mysql/users"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConfigDB struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string
	DB_HOST     string
	DB_PORT     string
}

func (config *ConfigDB) InitDB() *gorm.DB {
	var err error
	loc := "Asia%2FJakarta"

	var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("error when connecting to database: %s", err)
	}

	log.Println("connected to database")

	return db
}

func DBMigrate(db *gorm.DB) {
	db.AutoMigrate(&users.User{})
}

func CloseDB(db *gorm.DB) error {
	database, err := db.DB()

	if err != nil {
		log.Printf("error when getting database instance: %v", err)
		return err
	}

	if err := database.Close(); err != nil {
		log.Printf("error when closing database connection: %v", err)
		return err
	}

	log.Println("database connection is closed")

	return nil
}