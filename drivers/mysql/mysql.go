package mysql_driver

import (
	"fmt"
	"log"

	"backend/drivers/mysql/users"

	_utils "backend/utils"
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
		_utils.GetConfig("DB_USERNAME"),
		_utils.GetConfig("DB_PASSWORD"),
		_utils.GetConfig("DB_HOST"),
		_utils.GetConfig("DB_PORT"),
		_utils.GetConfig("DB_NAME"),
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