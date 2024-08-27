package template

import (
	"boilerplate/models"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	l := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

func TestSelect(t *testing.T) {
	db, mock := NewMockDB()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "John Doe")

	mock.ExpectQuery("^SELECT (.+) FROM `?users`?$").WillReturnRows(rows)

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		t.Fatalf("Error in finding users: %v", err)
	}

	if len(users) != 1 || users[0].Name != "John Doe" {
		t.Fatalf("Unexpected user data retrieved: %v", users)
	}
}

func TestInsert(t *testing.T) {
	db, mock := NewMockDB()

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `?users`? (.+)$").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := models.User{Name: "Jane Doe"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock := NewMockDB()

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `?users`? SET .*`?name`?=\\? WHERE (.+)$").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := models.User{
		Model: models.Model{ID: 1},
		Name:  "Jane Smith"}
	if err := db.Save(&user).Error; err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock := NewMockDB()

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM `?users`? WHERE (.+)$").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user := models.User{
		Model: models.Model{ID: 1},
	}
	if err := db.Delete(&user).Error; err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
}

func TestOneToMany(t *testing.T) {
	db, mock := NewMockDB()

	userRows := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow(1, "John Doe").
		AddRow(2, "Jane Doe")
	postRows := sqlmock.
		NewRows([]string{"id", "title", "user_id"}).
		AddRow(1, "Post 1", 1).
		AddRow(2, "Post 2", 1).
		AddRow(3, "Post 3", 2)

	mock.ExpectQuery("^SELECT (.+) FROM `?users`?$").WillReturnRows(userRows)
	mock.ExpectQuery("^SELECT (.+) FROM `?posts`? WHERE (.+)$").WillReturnRows(postRows)

	var users []models.User
	if err := db.Model(&models.User{}).Preload("Posts").Find(&users).Error; err != nil {
		t.Fatalf("Error in finding users: %v", err)
	}

	if len(users) != 2 || users[0].Name != "John Doe" || len(users[0].Posts) != 2 || users[0].Posts[0].Title != "Post 1" || users[1].Posts[0].UserID != 2 {
		t.Fatalf("Unexpected user data retrieved: %v", users)
	}
}
