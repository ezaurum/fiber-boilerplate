package models

// User model
type User struct {
	Model
	Name  string `json:"name"`
	Posts []Post `json:"posts"`
}
