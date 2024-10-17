package database

import (
	"boilerplate/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (c Connection) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		c.Host, c.User, c.Password, c.Name, c.Port,
	)
}

// Connect with database
func Connect(c Connection) *gorm.DB {
	_db, err := gorm.Open(postgres.Open(c.DSN()), &gorm.Config{})
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
	tx := db.Create(user)
	if tx.Error != nil {
		panic(tx.Error)
	}
}

func Get() []*models.User {
	var users []*models.User
	db.Find(&users)
	return users
}
