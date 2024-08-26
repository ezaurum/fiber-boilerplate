package models

// User model
type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Posts []Post `json:"posts"`
}
