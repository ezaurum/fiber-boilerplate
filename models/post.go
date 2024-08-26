package models

type Post struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	UserID uint   `json:"userID"`
}
