package models

import "time"

type Model struct {
	ID int64 `json:"id,string,omitempty" gorm:"primaryKey,autoIncrement=false"`
	UpdateModel
}

type UpdateModel struct {
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
