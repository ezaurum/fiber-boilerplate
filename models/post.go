package models

type Post struct {
	Model
	Title  string `json:"title"`
	UserID uint   `json:"userID"`
}
