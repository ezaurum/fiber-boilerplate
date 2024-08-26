package database

import (
	"boilerplate/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// Connect with database
func Connect() *gorm.DB {
	_db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = _db.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate database")
	}
	db = _db
	fmt.Println("Connected with Database")
	return db
}

func Insert(user *models.User) {
	db.Create(user)
}

func Get() []*models.User {
	var users []*models.User
	db.Find(&users)
	return users
}
